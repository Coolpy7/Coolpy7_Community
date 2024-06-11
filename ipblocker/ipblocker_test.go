package ipblocker_test

import (
	"bytes"
	"github.com/Coolpy7/Coolpy7_Community/ipblocker"
	"os/exec"
	"strings"
	"testing"
)

// 添加测试规则到iptables
func addTestRule(identifier string) error {
	comment := "blocker-test-" + identifier
	cmd := exec.Command("iptables", "-A", "INPUT", "-s", "192.168.0.1", "-m", "comment", "--comment", comment, "-j", "DROP")
	return cmd.Run()
}

// 检查特定规则是否存在
func ruleExists(ruleSpec string) bool {
	cmd := exec.Command("iptables", "-S")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run() // 这里我们忽略错误处理以简化示例
	return strings.Contains(out.String(), ruleSpec)
}

// 测试 cleanupOldRules 功能
func TestCleanupOldRules(t *testing.T) {
	identifier := "test123"
	b := ipblocker.NewBlocker(60, 3, identifier)
	defer b.Close()

	// 1. 添加一个测试规则
	if err := addTestRule(identifier); err != nil {
		t.Fatalf("Failed to add test rule: %v", err)
	}

	// 2. 确认规则已添加
	ruleSpec := "-A INPUT -s 192.168.0.1 -m comment --comment blocker-test-" + identifier + " -j DROP"
	if !ruleExists(ruleSpec) {
		t.Fatalf("Test rule not found after adding")
	}

	// 3. 调用 cleanupOldRules 清理规则
	b.CleanupOldRules()

	// 4. 确认规则已删除
	if ruleExists(ruleSpec) {
		t.Fatalf("Test rule found after cleanup")
	}
}
