package ipblocker

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type BlockedIP struct {
	unblockTime time.Time
}

// Blocker 用于管理IP封禁
type Blocker struct {
	failCounts  sync.Map
	blockedIPs  sync.Map
	BlockTime   time.Duration
	MaxAttempts int
	done        chan struct{}
	cleanerWG   sync.WaitGroup
	Identifier  string
}

// NewBlocker 创建一个新的Blocker
func NewBlocker(blockTime int, maxAttempts int, identifier string) *Blocker {
	b := &Blocker{
		BlockTime:   time.Duration(blockTime) * time.Second,
		MaxAttempts: maxAttempts,
		done:        make(chan struct{}),
		Identifier:  identifier,
	}

	// 清理旧的iptables规则
	b.CleanupOldRules()

	b.cleanerWG.Add(1)
	go b.cleanupTask() // 启动清理任务

	return b
}

// BlockIP 封禁IP
func (b *Blocker) BlockIP(ip string) error {
	if ip == "127.0.0.1" || ip == "::1" || ip == "localhost" { // 忽略本地地址
		return nil
	}

	comment := fmt.Sprintf("blocker-%s", b.Identifier)
	cmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-m", "comment", "--comment", comment, "-j", "DROP")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to block IP %s: %w", ip, err)
	}

	unblockTime := time.Now().Add(b.BlockTime)
	b.blockedIPs.Store(ip, BlockedIP{unblockTime: unblockTime})
	return nil
}

// UnblockIP 解封IP
func (b *Blocker) UnblockIP(ip string) error {
	// 使用iptables解封IP
	comment := fmt.Sprintf("blocker-%s", b.Identifier)
	cmd := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-m", "comment", "--comment", comment, "-j", "DROP")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to unblock IP %s: %w", ip, err)
	}

	b.blockedIPs.Delete(ip) // 从封禁的IP列表中删除
	return nil
}

// RegisterFailedAttempt 注册一次失败的尝试
func (b *Blocker) RegisterFailedAttempt(ip string) bool {
	if ip == "127.0.0.1" || ip == "::1" || ip == "localhost" { // 忽略本地地址
		return false
	}

	failCount, _ := b.failCounts.LoadOrStore(ip, 1)
	counter := 0 // default start counter

	switch v := failCount.(type) {
	case int:
		counter = v + 1 // increase counter
		b.failCounts.Store(ip, counter)
	default:
		counter = 1 // reset counter to 1
		b.failCounts.Store(ip, counter)
	}

	if counter >= b.MaxAttempts {
		_ = b.BlockIP(ip)
		b.failCounts.Delete(ip) // 重置计数器
		return true             // 表示已经封禁
	}

	return false // 未达到封禁次数
}

// cleanupTask 定期清理已过期的封禁
func (b *Blocker) cleanupTask() {
	defer b.cleanerWG.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			b.blockedIPs.Range(func(key, value interface{}) bool {
				ip := key.(string)
				blockedIP := value.(BlockedIP)
				if now.After(blockedIP.unblockTime) {
					// 无需处理错误，因为方法UnblockIP已经处理过
					_ = b.UnblockIP(ip)
				}
				return true
			})
		case <-b.done:
			return
		}
	}
}

// CleanupOldRules 清理旧的iptables规则
func (b *Blocker) CleanupOldRules() {
	cmd := exec.Command("iptables", "-S")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to list iptables rules:", err)
		return
	}

	commentPrefix := fmt.Sprintf("blocker-%s", b.Identifier)
	for _, line := range strings.Split(out.String(), "\n") {
		if strings.Contains(line, commentPrefix) {
			ruleSpec := strings.Fields(line)
			// 将规则删除
			delCmd := exec.Command("iptables", append([]string{"-D"}, ruleSpec[0:]...)...)
			if err := delCmd.Run(); err != nil {
				fmt.Printf("failed to delete iptables rule: %s, error: %v\n", line, err)
			} else {
				fmt.Printf("deleted iptables rule: %s\n", line)
			}
		}
	}
}

// cleanupIPTables 清除所有封禁
func (b *Blocker) cleanupIPTables() {
	b.blockedIPs.Range(func(key, value interface{}) bool {
		ip := key.(string)
		// 无需处理错误，因为方法UnblockIP已经处理过
		_ = b.UnblockIP(ip)
		return true
	})
}

// Close 关闭Blocker，释放资源并清除IP封禁规则
func (b *Blocker) Close() {
	close(b.done)
	b.cleanerWG.Wait()
	b.cleanupIPTables()
}
