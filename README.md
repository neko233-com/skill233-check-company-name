# skill233-check-company-name

> 公司名查重与子公司名储备 Agent Skill —— **全球范围，中国优先**。
> 用于「创建公司」前的重名/近似核验，以及围绕主品牌批量**储备衍生子公司名**。

## 能做什么

- **重名核验**：检查拟用公司名是否与现存企业登记名、商标重复或高度近似。
- **中国优先**：默认先查中国大陆（国家企业信用信息公示系统、中国商标网），再按需扩展到中国香港 / 中国澳门 / 中国台湾及海外主要法域。
- **可注册性初判**：按「完全重复 / 高度近似 / 低风险」分档，附检索依据与链接。
- **子公司名储备**：围绕主名批量生成候选（字号+行业、字号+地域、关系词、中英译名等），逐个查重，输出储备表。
- **保守可靠**：结论区分「已查证事实」与「初步判断」，始终附免责声明——不保证一定能注册。

## 适用人群

创业者、准备新设公司的团队、需要搭建多主体/子公司架构的企业、品牌与法务人员。

## 用法

这是一个 Agent Skill（见 [`SKILL.md`](./SKILL.md) 的 frontmatter）。在支持 Agent Skills 的客户端中安装后，直接用自然语言描述需求即可触发：

```
帮我核验「星澜科技」在中国大陆有没有重名，顺便储备 8 个子公司名
```

```
查一下「Aurora Labs」在美国和欧盟的商标与企业名重复情况
```

Skill 会主动询问缺失信息（拟用名 / 主营业务 / 优先法域 / 用途 / 储备数量），
然后按 `SKILL.md` 的工作流分级核验并输出结论表。

## 手动安装

```bash
git clone https://github.com/neko233-com/skill233-check-company-name.git
```

将仓库放入你的 skills 目录（如 Claude Code 的 `.claude/skills/`、Codex 的 skills 路径），
确保 `SKILL.md` 位于技能根目录。

## 目录结构

```
.
├── SKILL.md                       # Agent Skill 定义（触发条件 + 工作流）
├── README.md                      # 本文件
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
