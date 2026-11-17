# skill233-check-company-name

> 公司名查重与子公司名储备 Agent Skill —— **全球范围，中国优先**。
> 用于「创建公司」前的重名/近似核验，结构化告知**原名不可用原因**，并默认一次性生成 **5 个**衍生候选公司名。

## 能做什么

- **原名诊断**：判断拟用名是否与现存企业/商标重复或近似，并逐条给出**具体不可用原因**（重名 / 字号近似 / 商标冲突 / 需核准字样等），附依据链接。
- **5 个待选项**：默认一次性生成 5 个候选名（可用 `candidateCount` 指定数量），覆盖字号+行业同义词、字号+地域、关系词、英文译名、缩写+业务等策略，每个都带查重结论。
- **中国优先**：默认先查中国大陆（国家企业信用信息公示系统、中国商标网），再按需扩展到中国香港 / 中国澳门 / 中国台湾及海外主要法域。
- **三重核验**：企业名（工商登记）+ 商标（近似类别）+ 域名/社媒句柄，独立给出结论。
- **保守可靠**：结论区分「已查证事实」与「初步判断」，始终附免责声明——不保证一定能注册。

## 适用人群

创业者、准备新设公司的团队、需要搭建多主体/子公司架构的企业、品牌与法务人员。

## 技术栈

**Go 1.27**（标准库，零外部依赖）。判定引擎为纯函数、无网络 I/O，便于单元测试与嵌入 Agent Skill；联网检索由 Skill 工作流负责，检索到的事实以 JSON 喂给引擎。

```bash
go test ./...   # 单元测试（必须全绿）
go vet ./...
go build ./...
```

## 用法

### 作为 Agent Skill

在支持 Agent Skills 的客户端中安装后（确保 `SKILL.md` 在技能根目录），直接用自然语言触发：

```
帮我核验「星澜科技」在中国大陆有没有重名，告诉我为什么，并给 5 个备选名
```

```
查「Aurora Labs」在美国和欧盟的商标/企业名重复情况，再储备 5 个子公司名
```

Skill 会采集各渠道事实、调用 Go 引擎，输出「原因表 + 5 个待选项表 + 下一步 + 免责声明」。

### 作为命令行工具

```bash
echo '{"name":"星澜科技","industry":"软件开发","region":"上海"}' | go run ./cmd/checkname
```

输入 JSON 关键字段：

| 字段 | 说明 |
|---|---|
| `name` | 拟用名（必填） |
| `industry` | 主营业务，用于行业近似判断与生成候选 |
| `region` | 地域，如「上海」 |
| `jurisdiction` | 主法域，默认 `CN-MAINLAND` |
| `candidateCount` | 候选数量，默认 5 |
| `existingNames` | 已登记企业名数组（来自检索） |
| `existingTrademarks` | 已存在商标数组（来自检索） |

输出 JSON：`original`（风险+原因）、`unavailableReasons`（原名不可用原因）、`candidates`（默认 5 个）、`disclaimer`。

### 手动安装仓库

```bash
git clone https://github.com/neko233-com/skill233-check-company-name.git
```

将仓库放入你的 skills 目录（如 Claude Code 的 `.claude/skills/`、Codex 的 skills 路径），
确保 `SKILL.md` 位于技能根目录。

## 目录结构

```
.
├── SKILL.md                       # Agent Skill 定义（含「必须给原因 + 5 个候选」约束）
├── README.md                      # 本文件
├── go.mod                         # Go 1.27 模块（零外部依赖）
├── cmd/checkname/                 # CLI：stdin JSON -> stdout JSON
├── internal/namecheck/            # 判定引擎 + 单元测试
│   ├── types.go                   # 数据模型 / 常量
│   ├── normalize.go               # 名称规范化 / 字号提取 / 相似度
│   ├── assess.go                  # 风险与原因诊断
│   ├── suggest.go                 # 候选名生成
│   ├── check.go                   # 顶层流水线
│   └── *_test.go                  # 完整单元测试
└── references/
    └── jurisdictions.md           # 各法域官方查询渠道清单
```

## 查询渠道（速览）

| 法域 | 企业名 | 商标 |
|---|---|---|
| 中国大陆 | 国家企业信用信息公示系统 gsxt.gov.cn | 中国商标网 sbj.cnipa.gov.cn |
| 中国香港 | 公司注册处公司名称查册 | 知识产权署商标检索 |
| 中国澳门 | 商业及动产登记局 | 经济局知识产权 |
| 中国台湾 | 经济部商业司公司名称查询 | 智慧财产局商标检索 |
| 美国 | 各州 Secretary of State | USPTO TESS |
| 欧盟 | 各国公司登记处 | EUIPO eSearch |
| 全球 | OpenCorporates | WIPO Global Brand Database |

完整清单见 [`references/jurisdictions.md`](./references/jurisdictions.md)。

## 免责声明

本 Skill 提供基于公开渠道的初步查重与命名建议，不构成法律意见，
也不保证任何名称一定可注册。正式注册前请以登记机关与商标局的审查结果为准。

## License

MIT © neko233-com
