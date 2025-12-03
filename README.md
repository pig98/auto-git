# Auto Git Sync

一个 macOS 上的自动 Git 同步服务，监听配置的 Git 目录下的文件变动，自动同步到远程仓库。

## 功能特性

- 🔍 自动监听配置的 Git 目录下的文件变动
- 🔄 每次同步前自动拉取最新数据
- 🔀 自动尝试合并冲突
- ⚠️ 冲突无法解决时发送系统通知
- 🚀 使用 brew service 管理，开机自启
- 📝 支持日志级别配置（DEBUG/INFO/ERROR）
- 🔕 支持禁用系统通知

## 安装

### 使用 Homebrew（推荐）

1. **添加 tap**：
```bash
brew tap pig98/tap
```

2. **安装**：
```bash
brew install auto-git
```

**说明**：
- 安装时会从源码编译，Homebrew 会自动安装 Go 环境（如果尚未安装）
- 这是 Homebrew 的标准安装方式，与大多数 Formula 一致

3. **配置服务**：

首次启动前需要配置 `GIT_DIRS` 环境变量（其他配置已有默认值）：

```bash
# 编辑服务配置
brew services edit auto-git
```

在打开的 plist 文件中，找到 `EnvironmentVariables` → `GIT_DIRS`，修改为你的 Git 目录路径：

```xml
<key>EnvironmentVariables</key>
<dict>
    <key>GIT_DIRS</key>
    <string>/Users/yourname/project1:/Users/yourname/project2</string>
    <!-- 其他配置已有默认值，通常无需修改：
         QUIET_PERIOD_MINUTES: 10
         LOG_LEVEL: INFO
         DISABLE_NOTIFICATIONS: 0
    -->
</dict>
```

**配置说明**：
- `GIT_DIRS`：**必需**，多个目录用 `:` 分隔，使用绝对路径
- `QUIET_PERIOD_MINUTES`：默认 `10`（10 分钟无变化才同步）
- `LOG_LEVEL`：默认 `INFO`（可选：`DEBUG`、`INFO`、`ERROR`）
- `DISABLE_NOTIFICATIONS`：默认 `0`（启用通知，设为 `1` 禁用）

4. **启动服务**：
```bash
brew services start auto-git
```

### 从源码安装

1. **克隆仓库**：
```bash
git clone git@github.com:pig98/auto-git.git
cd auto-git
```

2. **编译**：
```bash
make build
# 或
go build -o auto-git .
```

3. **配置并安装服务**：

需要手动创建 plist 文件。参考 Homebrew 安装方式的配置，创建 `~/Library/LaunchAgents/cn.dev.sc.autogit.plist`：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>cn.dev.sc.autogit</string>
    <key>ProgramArguments</key>
    <array>
        <string>/path/to/auto-git</string>  <!-- 修改为实际的二进制路径 -->
    </array>
    <key>EnvironmentVariables</key>
    <dict>
        <key>GIT_DIRS</key>
        <string>/path/to/your/repo1:/path/to/your/repo2</string>  <!-- 修改为你的 Git 目录 -->
        <key>QUIET_PERIOD_MINUTES</key>
        <string>10</string>
        <key>LOG_LEVEL</key>
        <string>INFO</string>
        <key>DISABLE_NOTIFICATIONS</key>
        <string>0</string>
    </dict>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/auto-git.out.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/auto-git.err.log</string>
</dict>
</plist>
```

然后启动服务：
```bash
launchctl load ~/Library/LaunchAgents/cn.dev.sc.autogit.plist
```

## 服务管理

### 使用 Homebrew（推荐）

```bash
# 启动服务
brew services start auto-git

# 停止服务
brew services stop auto-git

# 重启服务
brew services restart auto-git

# 查看状态
brew services list | grep auto-git
```

### 手动管理

```bash
# 启动
launchctl load ~/Library/LaunchAgents/cn.dev.sc.autogit.plist

# 停止
launchctl unload ~/Library/LaunchAgents/cn.dev.sc.autogit.plist

# 查看状态
launchctl list | grep autogit
```

## 配置说明

### 环境变量

- `GIT_DIRS`: **必需**，要监听的 Git 仓库目录列表，多个目录用 `:` 分隔
  - 示例: `/Users/username/project1:/Users/username/project2`
- `QUIET_PERIOD_MINUTES`: **可选**，静默期时间（分钟），默认 10 分钟
  - 示例: `5` 表示 5 分钟内无文件修改才同步
  - 示例: `30` 表示 30 分钟内无文件修改才同步
- `LOG_LEVEL`: **可选**，日志级别，支持 `DEBUG` / `INFO` / `ERROR`，默认 `INFO`
  - `DEBUG`: 输出所有日志（包括频繁的文件变更日志）
  - `INFO`: 只输出关键流程日志（默认）
  - `ERROR`: 只输出错误日志
- `DISABLE_NOTIFICATIONS`: **可选**，是否禁用系统通知，默认启用通知
  - 设为 `1` / `true` / `yes` 时禁用通知

### plist 文件完整示例

运行 `brew services edit auto-git` 后，你会看到类似这样的配置：

```xml
<key>EnvironmentVariables</key>
<dict>
    <key>GIT_DIRS</key>
    <string></string>  <!-- ⚠️ 必须修改：填写你的 Git 目录路径，多个用 : 分隔 -->
    <key>QUIET_PERIOD_MINUTES</key>
    <string>10</string>  <!-- ✅ 默认值，通常无需修改 -->
    <key>LOG_LEVEL</key>
    <string>INFO</string>  <!-- ✅ 默认值，通常无需修改 -->
    <key>DISABLE_NOTIFICATIONS</key>
    <string>0</string>  <!-- ✅ 默认值，通常无需修改 -->
</dict>
```

**最小配置示例**（只需修改 GIT_DIRS）：
```xml
<key>EnvironmentVariables</key>
<dict>
    <key>GIT_DIRS</key>
    <string>/Users/yourname/project1:/Users/yourname/project2</string>
</dict>
```

其他配置项会使用 Formula 中定义的默认值。

## 工作原理

1. **文件监听**: 使用 `fsnotify` 库递归监听配置的 Git 目录
2. **静默期机制**: 每次有文件变化时重置定时器，只有在指定时间内（默认 10 分钟）完全没有新的文件变化时，才触发同步
3. **同步流程**:
   - 先暂存并提交本地变更
   - 拉取最新代码 (`git pull --rebase`)
   - 检查并解决冲突
   - 推送到远程 (`git push`)
4. **冲突处理**: 
   - 优先尝试保留本地更改（ours）
   - 如果失败，尝试使用远程更改（theirs）
   - 如果仍无法解决，发送系统通知
5. **系统通知**: 使用 macOS 的 `osascript` 发送通知

### 静默期机制说明

程序使用静默期机制来避免频繁同步：
- 每次检测到文件变化时，会重置定时器
- 只有在指定时间内（如 10 分钟）完全没有新的文件变化时，才会触发同步
- 这样可以避免在频繁编辑文件时产生大量提交，而是在你停止编辑一段时间后才进行同步
- 服务启动时，如果检测到已有未提交的变更，也会启动静默期定时器

## 开发

### 使用 Makefile

```bash
# 构建到 bin/ 目录
make build

# 清理编译产物
make clean

# 打包源码为 tar.gz（用于发布）
make package VERSION=0.1.0
```

### 查看版本信息

```bash
# 查看版本
auto-git -v
# 或
auto-git --version

# 输出示例：
# auto-git version 0.1.0
# Build time: 2025-12-03T09:47:09
# Git commit: abc1234
```

### 项目结构

```
auto-git/
├── main.go                    # 主程序入口
├── internal/                  # 内部包（不对外暴露）
│   ├── config/               # 配置管理
│   ├── git/                  # Git 操作
│   ├── logger/               # 日志管理
│   ├── notify/               # 系统通知
│   └── watcher/              # 文件监听
├── Formula/                  # Homebrew Formula 模板（参考用）
├── .github/                  # GitHub Actions workflows
│   └── workflows/           # CI/CD 配置
└── Makefile                  # 构建管理
```

### 日志位置

- Homebrew 安装：`$(brew --prefix)/var/log/auto-git.log` 和 `$(brew --prefix)/var/log/auto-git.err.log`
- 手动安装：根据 plist 中配置的 `StandardOutPath` 和 `StandardErrorPath`

## 故障排查

### 服务无法启动

1. 检查 `GIT_DIRS` 环境变量是否已配置
2. 查看错误日志
3. 检查 Homebrew 是否已安装: `brew --version`

### 同步失败

1. 检查 Git 凭证配置
2. 检查网络连接
3. 查看日志文件了解详细错误信息

### 权限问题

#### macOS Gatekeeper 阻止运行

如果安装后无法运行，提示"无法打开，因为无法验证开发者"，这是因为二进制文件未签名。解决方法：

**方法一：在系统设置中允许（推荐）**
1. 打开"系统设置" → "隐私与安全性"
2. 找到被阻止的应用，点击"仍要打开"
3. 或者在终端运行：
```bash
xattr -d com.apple.quarantine $(brew --prefix)/bin/auto-git
```

**方法二：确保可执行文件有执行权限**
```bash
chmod +x $(brew --prefix)/bin/auto-git
```

## 注意事项

- 确保配置的目录都是有效的 Git 仓库
- 确保有相应的 Git 权限（SSH 密钥或凭证已配置）
- 冲突解决策略可能不适合所有场景，建议定期检查
- 首次运行前建议手动测试一次同步流程
- 目录路径必须使用绝对路径
- 被 `.gitignore` 忽略的文件不会触发同步

## 技术细节

### brew service 配置流程

1. **Formula 中的 `service` 块**（`Formula/auto-git.rb`）：
   - 定义了服务的运行参数、日志路径、环境变量默认值等
   - 这是服务的"模板"配置
   - 默认配置包括：
     - `QUIET_PERIOD_MINUTES=10`
     - `LOG_LEVEL=INFO`
     - `DISABLE_NOTIFICATIONS=0`
     - `GIT_DIRS=""`（用户必须填写）

2. **安装时**（`brew install auto-git`）：
   - Homebrew 执行 `install` 方法编译和安装二进制
   - 此时不会生成 plist 文件

3. **首次启动服务时**（`brew services start auto-git`）：
   - Homebrew 读取 Formula 中的 `service` 块
   - 自动生成 plist 文件：`~/Library/LaunchAgents/homebrew.mxcl.auto-git.plist`
   - 使用 Formula 中的默认值填充环境变量
   - 使用 `launchctl load` 加载服务

4. **编辑配置**（`brew services edit auto-git`）：
   - 直接编辑已生成的 plist 文件
   - **用户只需修改 `GIT_DIRS`**，其他配置已有合理的默认值
   - 修改后需要重启服务：`brew services restart auto-git`

## 维护者指南

### 发布新版本

#### 方式一：使用 GitHub Actions（推荐）

1. **创建版本标签并推送**：
```bash
git tag v0.1.0
git push origin v0.1.0
```

**重新发布（如果 tag 已存在）**：

如果需要重新触发发布（比如修复了 workflow 或 Formula 问题），可以删除并重新创建 tag：

```bash
# 删除远程 tag
git push origin --delete v0.1.0

# 删除本地 tag（可选）
git tag -d v0.1.0

# 重新创建并推送 tag
git tag v0.1.0
git push origin v0.1.0
```

或者创建新的版本号：
```bash
git tag v0.1.1
git push origin v0.1.1
```

GitHub Actions 会自动：
- 编译并测试程序
- 等待 GitHub 生成源码包（GitHub 在创建 tag 时会自动生成）
- 计算源码包的 SHA256
- 更新 homebrew-tap 仓库中的 Formula
- 提交并推送到 homebrew-tap

**注意**：需要在 GitHub 仓库设置中添加 Secret：

##### 步骤 1：创建 Personal Access Token

1. 登录 GitHub，点击右上角头像 → **Settings**
2. 左侧菜单滚动到底部，点击 **Developer settings**
3. 点击 **Personal access tokens** → **Tokens (classic)**
4. 点击 **Generate new token** → **Generate new token (classic)**
5. 填写信息：
   - **Note**（备注）：`Homebrew Tap Token`（或任意描述）
   - **Expiration**（过期时间）：选择合适的时间（建议 90 天或更长）
   - **Select scopes**（权限范围）：勾选 **`repo`**（完整仓库权限）
     - 这会自动勾选所有 repo 相关权限，包括读写权限
6. 点击 **Generate token**
7. **重要**：立即复制生成的 token（只显示一次，关闭后无法再查看）

##### 步骤 2：添加 Secret 到仓库

1. 进入你的源码仓库：`https://github.com/pig98/auto-git`
2. 点击仓库顶部的 **Settings** 标签
3. 左侧菜单点击 **Secrets and variables** → **Actions**
4. 点击 **New repository secret**
5. 填写信息：
   - **Name**：`HOMEBREW_TAP_TOKEN`
   - **Secret**：粘贴刚才复制的 token
6. 点击 **Add secret**

完成！现在 GitHub Actions 可以使用这个 token 来推送 homebrew-tap 仓库了。

#### 验证发布流程

发布后，按以下步骤验证：

1. **检查 GitHub Actions 是否成功**：
   - 访问 `https://github.com/pig98/auto-git/actions`
   - 确认 Release workflow 成功完成

2. **验证源码包是否生成**：
```bash
# 替换为实际版本号
VERSION="0.1.0"
curl -I "https://github.com/pig98/auto-git/archive/refs/tags/v${VERSION}.tar.gz"
# 应该返回 200 OK
```

3. **验证 Formula 是否更新**：
```bash
# 检查 homebrew-tap 仓库
git clone git@github.com:pig98/homebrew-tap.git /tmp/homebrew-tap-check
cd /tmp/homebrew-tap-check
cat Formula/auto-git.rb | grep -E "(version|url|sha256)"
# 应该看到正确的版本、URL 和 SHA256
```

4. **测试安装**：
```bash
# 添加 tap
brew tap pig98/tap

# 尝试安装（使用 --dry-run 先测试）
brew install --dry-run auto-git

# 如果测试通过，正式安装
brew install auto-git
```

5. **验证安装**：
```bash
# 检查二进制是否存在
which auto-git
auto-git -v

# 检查服务配置
brew services edit auto-git
```

#### 清理 Homebrew 缓存（如果重新安装后版本信息不对）

如果修改了 Formula 后重新安装，但版本信息还是旧的，可能是 Homebrew 缓存问题：

```bash
# 1. 更新 tap（获取最新的 Formula）
brew update

# 2. 卸载旧版本
brew uninstall auto-git

# 3. 清理缓存和下载的源码包
brew cleanup auto-git
rm -rf $(brew --cache)/auto-git-*

# 4. 强制重新安装（从源码编译）
brew install --force --build-from-source auto-git

# 5. 验证版本信息
auto-git -v
```

或者更彻底的清理方式：
```bash
# 完全清理并重新安装
brew uninstall auto-git
brew cleanup -s auto-git
brew update
brew install --build-from-source auto-git
auto-git -v
```

### Formula 管理说明

- **源码仓库中的 `Formula/` 目录**：仅作为模板/参考，不会被 Homebrew 使用
- **实际的 Formula 文件**：位于 `homebrew-tap` 仓库的 `Formula/auto-git.rb`
- **自动同步**：GitHub Actions 在发布新版本时自动更新 homebrew-tap 中的 Formula

## License

MIT License - see [LICENSE](LICENSE) file for details

## 贡献

欢迎提交 Issue 和 Pull Request！
