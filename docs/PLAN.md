# Aequa-network 基础 DVT 引擎第一阶段执行计划（SST: docs/Aequa-network.pdf）

本计划是第一阶段开发的唯一执行依据（SST）。除架构与功能外，安全与测试等同优先级，所有合并均以“审计就绪”为唯一衡量标准。

## 强制原则（安全铁律与审计就绪）
- 安全嵌入日常：从第一行代码起落实 SDL——STRIDE、OWASP Go、govulncheck+Snyk、安全 CR 清单、分支保护，接入 CI/PR 门禁。
- 测试同等优先：Fuzz + 故障注入与功能并行；PR 跑冒烟，夜间/周任务跑全量。
- 唯一标准：是否能通过顶级安全审计。

## 开发要求（方法论）
- 第一步：以我们的 PDF 为纲（SST）。
- 第二步：独立完成初始设计与编码。
- 第三步：外部实现仅作“顾问”，用于交叉验证复杂通用难题。
- 第四步：以我们架构为准绳（尤其“通用负载接口”）。

## 时间线（6–8 周，可并行）
- M0：工程与 SDL 基线（已完成）
- M1：核心服务骨架 + ValidatorAPI（进行中）
- M2：P2P（零信任）+ Gating + 资源管理 + DKG/lock
- M3：QBFT + StateDB + 悲观恢复 + PayloadManager
- M4：Fuzz + 故障注入/混沌测试框架 + 指标与告警

### M0 工程与 SDL 基线（Owner: Core）
- [x] 分层脚手架、Makefile、Dockerfile、docker-compose：一键 4 节点（已验 /health、/metrics 200）
- [x] CI 门禁：golangci-lint、单测、覆盖率阈值（当前0%，逐步升≥90%）、govulncheck、Snyk、Go vet、分支保护
- [x] SDL 文档与流程（威胁建模骨架、OWASP Go、安全 Code Review 清单）
- [x] 可观测性规范（JSON/Prom/OTel 骨架）

DoD：CI 全绿；Fuzz/Chaos 冒烟接入；本地一键起停；SDL/可观测性文档落地。

### M1 核心服务与 ValidatorAPI（Owner: Core）
- [x] 事件模型与总线（背压策略）
- [~] 生命周期扩展（错误传播/关停）
- [~] ValidatorAPI：代理非关键请求、拦截关键职责、严格校验；Fuzz 目标

### M2 P2P 与 DKG（Owner: Networking+Crypto）
- [ ] P2P + resource-manager
- [ ] Connection Gating 白名单与限流/评分断联
- [ ] DKG + cluster-lock 多签 + 启动硬门禁

### M3 QBFT 与持久化（Owner: Consensus+State）
- [ ] QBFT（EEA 规范）
- [ ] 严格消息验证与防重放
- [ ] StateDB 原子持久化与悲观恢复
- [ ] PayloadManager 标准负载

### M4 高级测试框架（Owner: QA/SRE）
- [ ] Fuzz 基建与字典
- [ ] 故障注入/混沌测试框架
- [ ] 安全指标与告警

## 并行测试与门禁
- PR 门禁：lint、unit（≥门槛）、govulncheck、Snyk、Fuzz/Chaos 冒烟
- 夜间：Fuzz 全量 + E2E 全量；周：Chaos 全量 + 压测

## 当前进度
- 总体：50%
- 已完成：M0
- 进行中：M1-1（事件总线与 API 骨架）

