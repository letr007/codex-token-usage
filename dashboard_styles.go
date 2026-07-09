package main

const dashboardStyles = `
:root{
  color-scheme:light dark;
  --bg:#faf9f5;
  --surface:#faf9f5;
  --surface-2:#f0eee8;
  --surface-3:#e9e6df;
  --text:#2d2a26;
  --muted:#6d6760;
  --line:#e3e1db;
  --line-strong:#d2cdc4;
  --row:#f7f5f0;
  --hover:rgba(45,42,38,.045);
  --blue:var(--cpa-primary,#8b8680);
  --cyan:var(--cpa-info,#7f8a86);
  --green:var(--cpa-success,#74806a);
  --orange:var(--cpa-warning,#a28461);
  --red:var(--cpa-danger,#a26761);
  --violet:var(--cpa-accent,#847971);
  --shadow:0 18px 42px rgba(45,42,38,.08);
  --shadow-soft:0 8px 24px rgba(45,42,38,.045);
  --panel-shadow:0 2px 10px rgba(45,42,38,.035);
  --focus:0 0 0 3px color-mix(in srgb,var(--blue) 14%,transparent);
}

:root[data-host-theme="white"],:root[data-host-theme="light"],:root[data-theme="white"],:root[data-theme="light"]{
  color-scheme:light;
  --bg:#ffffff;
  --surface:#ffffff;
  --surface-2:#ffffff;
  --surface-3:#f6f6f6;
  --text:#2d2a26;
  --muted:#6d6760;
  --line:#e5e5e5;
  --line-strong:#dddddd;
  --row:#fafafa;
  --hover:rgba(45,42,38,.038);
  --blue:var(--cpa-primary,#2d2a26);
  --cyan:var(--cpa-info,#7f8a86);
  --green:var(--cpa-success,#74806a);
  --orange:var(--cpa-warning,#a28461);
  --red:var(--cpa-danger,#a26761);
  --violet:var(--cpa-accent,#8b8680);
  --shadow:0 18px 38px rgba(20,18,16,.07);
  --shadow-soft:0 8px 24px rgba(20,18,16,.04);
  --panel-shadow:0 1px 8px rgba(20,18,16,.03);
  --focus:0 0 0 3px rgba(45,42,38,.10);
}

:root[data-host-theme="dark"],:root[data-theme="dark"]{
  color-scheme:dark;
  --bg:#151412;
  --surface:#1d1b18;
  --surface-2:#262320;
  --surface-3:#312c27;
  --text:#f6f4f1;
  --muted:#c9c3bb;
  --line:#3a3530;
  --line-strong:#45403a;
  --row:rgba(255,248,238,.028);
  --hover:rgba(240,232,222,.075);
  --blue:var(--cpa-primary,#8b8680);
  --cyan:var(--cpa-info,#8f9a94);
  --green:var(--cpa-success,#87947e);
  --orange:var(--cpa-warning,#af926b);
  --red:var(--cpa-danger,#b07872);
  --violet:var(--cpa-accent,#9c8d84);
  --shadow:0 22px 48px rgba(0,0,0,.34);
  --shadow-soft:0 12px 28px rgba(0,0,0,.2);
  --panel-shadow:0 2px 14px rgba(0,0,0,.18);
  --focus:0 0 0 3px color-mix(in srgb,var(--blue) 22%,transparent);
}

*{box-sizing:border-box}

body{
  margin:0;
  background:var(--bg);
  color:var(--text);
  font:13px/1.5 "SF Pro Text","Helvetica Neue","PingFang SC","Noto Sans SC",ui-sans-serif,system-ui,-apple-system,"Segoe UI",sans-serif;
}

main{max-width:1720px;margin:0 auto;padding:22px 24px 48px;--host-action-reserve:220px}

.hero{
  display:flex;
  justify-content:space-between;
  gap:16px;
  align-items:center;
  background:transparent;
  border:0;
  border-radius:18px;
  padding:4px calc(4px + var(--host-action-reserve)) 14px 2px;
  box-shadow:none;
  position:relative;
  top:auto;
  z-index:auto;
  backdrop-filter:none;
}

h1{
  font:700 24px/1.12 "SF Pro Text","Helvetica Neue","PingFang SC","Noto Sans SC",ui-sans-serif,system-ui,-apple-system,"Segoe UI",sans-serif;
  margin:0;
  letter-spacing:-.02em;
}

.hint{color:var(--muted);font-size:12px;line-height:1.5}
.controls{display:flex;gap:8px;align-items:center;flex-wrap:wrap;justify-content:flex-end}

input,select,button{
  height:34px;
  border:1px solid var(--line);
  background:var(--surface);
  color:var(--text);
  border-radius:11px;
  padding:0 11px;
  outline:none;
  transition:border-color .18s ease,background .18s ease,box-shadow .18s ease,filter .18s ease,transform .18s ease,color .18s ease;
}

input:focus-visible,select:focus-visible,button:focus-visible{box-shadow:var(--focus);border-color:var(--blue)}
input::placeholder{color:color-mix(in srgb,var(--muted) 84%,transparent)}
select,input{box-shadow:none}

input[type="checkbox"]{
  appearance:auto;
  width:13px;
  height:13px;
  min-width:13px;
  margin:0;
  padding:0;
  border:0;
  border-radius:3px;
  background:transparent;
  vertical-align:-2px;
  accent-color:var(--blue);
}

input[type="checkbox"]:focus-visible{box-shadow:var(--focus);outline:none}

button{
  background:var(--blue);
  border-color:var(--blue);
  color:var(--surface);
  cursor:pointer;
  font-weight:760;
  box-shadow:0 1px 2px rgba(0,0,0,.06),0 8px 18px rgba(45,42,38,.08);
}

button:hover{filter:brightness(.96);transform:translateY(-1px)}
button:disabled{opacity:.45;cursor:not-allowed;transform:none}

.ghost{
  background:var(--surface-3);
  color:var(--text);
  border-color:var(--line);
  box-shadow:none;
}

.danger-ghost{color:var(--red);border-color:color-mix(in srgb,var(--red) 28%,var(--line));background:color-mix(in srgb,var(--red) 6%,var(--surface-3))}

.modal-backdrop{position:fixed;inset:0;z-index:80;display:grid;place-items:center;padding:24px;background:rgba(21,20,18,.28);backdrop-filter:blur(8px)}
.modal-backdrop[hidden]{display:none}

.modal-panel{
  width:min(520px,100%);
  background:var(--surface);
  border:1px solid var(--line);
  border-radius:20px;
  box-shadow:var(--shadow);
  overflow:hidden;
  backdrop-filter:none;
}

.modal-head{display:flex;align-items:center;justify-content:space-between;gap:10px;padding:18px 20px 14px;border-bottom:1px solid var(--line);background:var(--surface)}
.modal-title-actions{display:flex;align-items:center;gap:8px;min-width:0;flex-wrap:wrap}
.modal-head h2,.section h2 span:first-child{font-family:"SF Pro Text","Helvetica Neue","PingFang SC","Noto Sans SC",ui-sans-serif,system-ui,-apple-system,"Segoe UI",sans-serif}
.modal-head h2{margin:0;font-size:16px;font-weight:700;letter-spacing:-.01em}
.compact-danger{height:28px;padding:0 10px;font-size:11px}
.icon-button{width:34px;min-width:34px;padding:0;font-size:18px;line-height:1}
.modal-body{display:grid;gap:12px;padding:18px 20px}
.form-row{display:grid;grid-template-columns:92px minmax(0,1fr);gap:10px;align-items:center}
.form-row span{color:var(--muted);font-size:12px;font-weight:700}
.form-row input{width:100%;height:36px}
.modal-note{border:1px solid var(--line);border-radius:14px;background:var(--surface-3);padding:10px 12px;color:var(--muted);font-size:12px;line-height:1.5;min-width:0}
.modal-note code{color:var(--text);font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.oauth-link-row{display:flex;align-items:center;gap:8px;min-width:0;flex-wrap:wrap}
.oauth-open-link{display:inline-flex;align-items:center;justify-content:center;min-height:30px;border:1px solid var(--line);border-radius:11px;background:var(--surface-3);padding:0 10px;color:var(--text);font-weight:760;text-decoration:none}
.oauth-copy-link{height:30px;font-size:11px}
.oauth-link-row code{max-width:260px;display:inline-block}
.modal-status{min-height:20px;color:var(--muted);font-size:12px;line-height:1.5}
.modal-status.ok{color:var(--green);font-weight:780}
.modal-status.warn{color:var(--orange);font-weight:780}
.modal-status.bad{color:var(--red);font-weight:780}
.modal-progress-status{display:flex;align-items:center;justify-content:space-between;gap:8px}
.modal-progress-status b{font-size:11px;font-variant-numeric:tabular-nums}
.modal-progress{height:6px;margin-top:4px;border:0;border-radius:999px;background:var(--surface-3);overflow:hidden}
.modal-progress span{display:block;height:100%;min-width:2px;border-radius:999px;background:var(--blue);transition:width .28s ease}
.modal-actions{display:flex;justify-content:flex-end;gap:8px;padding:14px 20px 18px;border-top:1px solid var(--line);background:var(--surface);flex-wrap:wrap}

.status{min-height:18px;margin:8px 2px;color:var(--muted);font-size:12px}
.tabs{display:flex;gap:8px;margin:0 0 16px;align-items:center;flex-wrap:wrap}
.tab-strip{display:flex;gap:6px;align-items:center;flex:1 1 auto;min-width:0;overflow:auto;scrollbar-width:thin;padding:4px;border:1px solid var(--line);border-radius:14px;background:var(--surface-3)}

.tab{
  height:34px;
  background:transparent;
  color:var(--muted);
  border-color:transparent;
  font-weight:700;
  white-space:nowrap;
  box-shadow:none;
}

.tab.active{
  background:var(--blue);
  border-color:var(--blue);
  color:var(--surface);
  box-shadow:0 1px 2px rgba(0,0,0,.06);
}

.tab-count{font-size:10.5px;margin-left:5px;opacity:.78}
.provider-picker{position:relative;margin-left:auto}
.picker-panel{display:none;position:absolute;right:0;top:42px;z-index:20;width:290px;max-height:360px;overflow:auto;background:var(--surface);border:1px solid var(--line);border-radius:16px;box-shadow:var(--shadow);padding:8px}
.provider-picker.open .picker-panel{display:block}
.picker-row{display:grid;grid-template-columns:18px minmax(0,1fr) max-content;gap:8px;align-items:center;min-height:34px;padding:5px 6px;border-radius:10px;color:var(--text)}
.picker-row:hover{background:var(--hover)}
.picker-row input{height:auto}
.picker-name{overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-weight:730}
.picker-meta{color:var(--muted);font-size:11px;font-variant-numeric:tabular-nums}
.fallback-key{display:none}
.fallback-key.on{display:inline-block}
[data-page]{display:none}
[data-page].page-on{display:block}
.command-grid{display:grid;grid-template-columns:minmax(0,1.75fr) minmax(320px,.72fr);gap:16px}
.cards{display:grid;grid-template-columns:repeat(5,minmax(0,1fr));gap:12px}

.metric{
  background:var(--surface);
  border:1px solid var(--line);
  border-radius:18px;
  padding:14px 16px;
  min-width:0;
  box-shadow:var(--panel-shadow);
  position:relative;
}

.metric:before{content:none}

.label{color:var(--muted);font-size:11px;white-space:nowrap;letter-spacing:.01em}
.label:before{content:"";display:inline-block;width:6px;height:6px;border-radius:999px;background:color-mix(in srgb,var(--accent,var(--cyan)) 74%,#fff 26%);box-shadow:none;margin-right:6px;vertical-align:1px}
.value{font-size:20px;font-weight:800;margin:4px 0 0;font-variant-numeric:tabular-nums;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;letter-spacing:-.02em}
.sub{font-size:11px;color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;margin-top:4px}

.layout{display:grid;grid-template-columns:1.52fr .78fr;gap:16px;margin-top:16px}

.section{
  background:var(--surface);
  border:1px solid var(--line);
  border-radius:20px;
  overflow:hidden;
  box-shadow:var(--panel-shadow);
}

.section h2{font-size:14px;margin:0;padding:18px 20px 12px;border-bottom:1px solid var(--line);display:flex;justify-content:space-between;gap:8px;background:var(--surface);letter-spacing:-.01em}
.section-body{padding:16px 20px 20px}
.scroll{overflow:auto;scrollbar-color:var(--line-strong) transparent;scrollbar-width:thin}
.insights{display:grid;grid-template-columns:1fr;gap:10px}
.insight{border:1px solid var(--line);background:var(--surface-3);border-radius:16px;padding:10px 12px;min-width:0;display:grid;grid-template-columns:84px minmax(0,1fr);column-gap:10px;align-items:center}
.insight b{display:block;font-size:13px;margin:0;font-variant-numeric:tabular-nums;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.insight span{display:block;color:var(--muted);font-size:11px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.insight span:first-child{grid-row:1 / 3;color:var(--text);font-weight:760}
.tone-red{border-color:color-mix(in srgb,var(--red) 34%,var(--line));background:color-mix(in srgb,var(--red) 7%,var(--surface-3))}
.tone-orange{border-color:color-mix(in srgb,var(--orange) 36%,var(--line));background:color-mix(in srgb,var(--orange) 7%,var(--surface-3))}
.tone-green{border-color:color-mix(in srgb,var(--green) 34%,var(--line));background:color-mix(in srgb,var(--green) 6%,var(--surface-3))}

.chart{width:100%;height:184px;display:block}
.trend-hit{fill:transparent;pointer-events:all;cursor:crosshair}
.trend-guide{stroke:var(--line-strong);stroke-width:1.2;stroke-dasharray:4 4}
.trend-dot{stroke:var(--surface);stroke-width:2}
.trend-tip-box{fill:var(--surface);stroke:var(--line);filter:drop-shadow(0 10px 18px rgba(28,24,20,.08))}
.trend-tip-title{fill:var(--text);font-size:12px;font-weight:800}
.trend-tip-line{fill:var(--muted);font-size:11px}
.legend{display:flex;gap:12px;align-items:center;justify-content:center;color:var(--muted);font-size:11px;margin-top:6px}
.dot{width:7px;height:7px;border-radius:50%;display:inline-block;margin-right:4px}
.mix{display:grid;gap:8px}
.mix-row{display:grid;grid-template-columns:72px 1fr 68px;gap:10px;align-items:center;font-size:11.5px}
.bar{height:8px;min-width:68px;background:var(--surface-3);border:0;border-radius:999px;overflow:hidden}
.bar span{display:block;min-width:2px;height:100%;background:var(--color,var(--blue));border-radius:999px}
.cell-meter{display:grid;gap:3px;min-width:108px}
.cell-meter b{font-size:12px;font-weight:760}
.cell-meter .bar{height:6px}
.pill{font-size:11px;border:1px solid var(--line);border-radius:999px;padding:3px 8px;color:var(--muted);white-space:nowrap;background:var(--surface-3)}
.quota-compact{display:grid;grid-template-columns:28px minmax(88px,1fr) minmax(126px,max-content);gap:8px;align-items:center;min-width:248px;font-size:11px;color:var(--muted);text-align:left}
.quota-compact+.quota-compact{margin-top:3px}
.quota-compact .bar{height:7px}
.danger{color:var(--red);font-weight:740}
.warn{color:var(--orange);font-weight:740}
.ok{color:var(--green);font-weight:740}

.account-toolbar{display:grid;grid-template-columns:minmax(220px,300px) 128px max-content minmax(12px,1fr) 86px 62px 54px 62px;gap:8px;align-items:center;margin-bottom:12px}
.account-toolbar>input{min-width:0;width:100%}
.account-toolbar select,.account-toolbar button{width:100%}
.account-toolbar .spacer{display:block}
.column-controls{display:flex;gap:6px;align-items:center;flex-wrap:nowrap;min-width:max-content}
.column-controls label{display:inline-flex;align-items:center;gap:5px;height:28px;border:1px solid var(--line);border-radius:10px;background:var(--surface-3);padding:0 8px;color:var(--muted);font-size:11px;font-weight:760;white-space:nowrap}
.column-controls label:hover{border-color:var(--line-strong);color:var(--text);background:var(--hover)}
.column-controls input{flex:0 0 auto}
.hide-account-perf [data-col="perf"],.hide-account-cache [data-col="cache"],.hide-account-quota5h [data-col="quota5h"],.hide-account-status [data-col="status"]{display:none}
.page-label{min-width:54px;text-align:center;color:var(--muted);font-size:12px;font-weight:750}
.account-summary-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(96px,1fr));gap:10px;margin-bottom:12px}
.account-summary-card{border:1px solid var(--line);border-radius:16px;background:var(--surface);padding:12px 14px;min-width:0;text-align:left;box-shadow:var(--panel-shadow)}
.account-summary-card span{display:block;color:var(--muted);font-size:11px}
.account-summary-card b{display:block;margin-top:2px;font-size:15px;font-variant-numeric:tabular-nums;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.account-summary-action{height:auto;color:var(--text);font:inherit;box-shadow:none;cursor:pointer;position:relative;overflow:hidden;transition:border-color .18s ease,background .18s ease,box-shadow .18s ease,transform .18s ease}
.account-summary-action:after{content:none}
.account-summary-action:hover{transform:translateY(-1px);box-shadow:var(--shadow-soft)}
.account-summary-action small{display:block;margin-top:1px;color:var(--muted);font-size:10px;font-weight:780}
.invalid-auth-action,.workspace-deactivated-action,.autoban-release-action{background:var(--surface-3);border-color:color-mix(in srgb,var(--line) 80%,var(--blue) 20%);backdrop-filter:none}
.invalid-auth-action.has-invalid{background:color-mix(in srgb,var(--red) 7%,var(--surface));border-color:color-mix(in srgb,var(--red) 44%,var(--line));box-shadow:none}
.workspace-deactivated-action.has-invalid{background:color-mix(in srgb,var(--orange) 8%,var(--surface));border-color:color-mix(in srgb,var(--orange) 46%,var(--line));box-shadow:none}
.autoban-release-action.has-invalid{background:color-mix(in srgb,var(--red) 8%,var(--surface));border-color:color-mix(in srgb,var(--red) 46%,var(--line));box-shadow:none}
.invalid-auth-action.has-invalid span,.invalid-auth-action.has-invalid small,.autoban-release-action.has-invalid span,.autoban-release-action.has-invalid small{color:color-mix(in srgb,var(--red) 76%,var(--text))}
.workspace-deactivated-action.has-invalid span,.workspace-deactivated-action.has-invalid small{color:color-mix(in srgb,var(--orange) 76%,var(--text))}
.invalid-auth-action.has-invalid b,.autoban-release-action.has-invalid b{color:var(--red)}
.workspace-deactivated-action.has-invalid b{color:var(--orange)}
.invalid-auth-panel{width:min(760px,100%)}
.workspace-deactivated-panel{width:min(720px,100%)}
.autoban-release-panel{width:min(760px,100%)}
.invalid-auth-toolbar{display:flex;align-items:center;justify-content:space-between;gap:8px}
.invalid-auth-toolbar span{color:var(--muted);font-size:12px;font-weight:780}
.invalid-auth-list{display:grid;gap:8px;max-height:420px;overflow:auto;padding-right:2px;scrollbar-width:thin}
.invalid-auth-row{display:grid;grid-template-columns:24px minmax(170px,1.15fr) minmax(128px,.78fr) minmax(120px,.8fr) 78px 60px;gap:8px;align-items:center;border:1px solid var(--line);border-radius:16px;background:var(--surface);padding:10px 12px;min-width:0;box-shadow:var(--panel-shadow)}
.workspace-deactivated-row{grid-template-columns:24px minmax(170px,1.15fr) minmax(128px,.78fr) minmax(120px,.8fr) 60px}
.autoban-release-row{grid-template-columns:24px minmax(170px,1.1fr) minmax(128px,.78fr) minmax(160px,.95fr) 68px}
.invalid-auth-row:hover{background:color-mix(in srgb,var(--red) 4%,var(--surface-2));border-color:color-mix(in srgb,var(--red) 18%,var(--line))}
.workspace-deactivated-row:hover{background:color-mix(in srgb,var(--orange) 5%,var(--surface-2));border-color:color-mix(in srgb,var(--orange) 20%,var(--line))}
.autoban-release-row:hover{background:color-mix(in srgb,var(--red) 5%,var(--surface-2));border-color:color-mix(in srgb,var(--red) 20%,var(--line))}
.invalid-auth-row.busy{border-color:color-mix(in srgb,var(--blue) 42%,var(--line));box-shadow:inset 3px 0 0 var(--blue)}
.invalid-auth-check{display:flex;align-items:center;justify-content:center}
.invalid-auth-main,.invalid-auth-meta{display:grid;gap:2px;min-width:0}
.invalid-auth-main b,.invalid-auth-main span,.invalid-auth-meta span,.invalid-auth-reason,.invalid-auth-reason b,.invalid-auth-reason span{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.invalid-auth-main b{font-size:12px;color:var(--text)}
.invalid-auth-main span,.invalid-auth-meta span,.invalid-auth-reason{font-size:11px;color:var(--muted)}
.invalid-auth-reason{color:var(--red);font-weight:760}
.invalid-auth-reason b{display:block;color:var(--red);font-size:11.5px}
.invalid-auth-reason span{display:block;color:var(--muted);font-size:11px;font-weight:650}
.workspace-deactivated-row .invalid-auth-reason{color:var(--orange)}
.workspace-deactivated-row .invalid-auth-reason b{color:var(--orange)}
.invalid-auth-pager{display:flex;align-items:center;justify-content:flex-end;gap:6px}
.invalid-auth-empty{border:1px dashed var(--line-strong);border-radius:16px;padding:18px;color:var(--muted);text-align:center;background:var(--surface-3)}
.autoban-toolbar{display:grid;grid-template-columns:minmax(180px,1fr) 86px 62px 54px 62px;gap:8px;align-items:center;margin-bottom:12px}
.autoban-toolbar select,.autoban-toolbar button{width:100%}
.autoban-scope{color:var(--muted);font-size:12px;font-weight:740;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.account-table-wrap{max-height:700px;border:1px solid var(--line);border-radius:18px;background:var(--surface);box-shadow:var(--panel-shadow)}
.autoban-table-wrap{max-height:310px;border:1px solid var(--line);border-radius:18px;background:var(--surface);box-shadow:var(--panel-shadow)}
.model-table-wrap{max-height:230px}
.recent-table-wrap{max-height:360px}
.account-table th,.account-table td{padding:10px 12px}
.account-cell{min-width:230px;max-width:360px}
.account-name{display:block;color:var(--text);font-weight:800;line-height:1.22;overflow-wrap:anywhere;white-space:normal}
.account-meta{display:flex;gap:6px;align-items:center;flex-wrap:wrap;margin-top:5px}
.account-id{color:var(--muted);font-size:11px;font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;overflow-wrap:anywhere;white-space:normal}
.metric-stack{display:grid;gap:2px;font-variant-numeric:tabular-nums}
.metric-stack b{font-size:12px}
.metric-stack span{color:var(--muted);font-size:11px}
.metric-stack .cost-line{color:var(--blue);font-weight:760}
.cost-weak{color:var(--muted);font-size:11px}
.recent-model{min-width:210px;max-width:360px}
.recent-primary{display:flex;gap:6px;align-items:center;min-width:0}
.model-chip{display:inline-flex;align-items:center;max-width:230px;min-height:24px;border:1px solid var(--line);border-radius:999px;background:var(--surface-3);padding:2px 9px;color:var(--text);font-weight:850;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.model-chip:before{content:"";width:7px;height:7px;border-radius:50%;background:var(--blue);margin-right:6px;flex:0 0 auto}
.recent-sub{display:block;margin-top:4px;color:var(--muted);font-size:11px;max-width:340px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.recent-badges{display:flex;gap:6px;align-items:center;justify-content:flex-start;flex-wrap:wrap}
.latency-pill,.cost-pill{display:inline-flex;align-items:center;justify-content:center;min-height:24px;border:1px solid var(--line);border-radius:999px;background:var(--surface-3);padding:2px 9px;font-weight:850;font-variant-numeric:tabular-nums;color:var(--text)}
.latency-pill.fast{color:var(--green);border-color:color-mix(in srgb,var(--green) 28%,var(--line));background:color-mix(in srgb,var(--green) 6%,var(--surface))}
.latency-pill.slow{color:var(--orange);border-color:color-mix(in srgb,var(--orange) 30%,var(--line));background:color-mix(in srgb,var(--orange) 7%,var(--surface))}
.cost-pill{background:color-mix(in srgb,var(--text) 4%,var(--surface));border-color:var(--line-strong)}
.token-main{font-weight:850;font-variant-numeric:tabular-nums}
.token-sub,.detail-sub{display:block;margin-top:3px;color:var(--muted);font-size:11px;font-variant-numeric:tabular-nums}
.detail-main{font-weight:780;color:var(--text)}
.status-pill{display:inline-flex;align-items:center;justify-content:center;min-height:22px;border:1px solid var(--line);border-radius:999px;padding:0 8px;background:var(--surface-3);font-size:11px;font-weight:800}
.status-pill.danger{border-color:color-mix(in srgb,var(--red) 34%,var(--line));background:color-mix(in srgb,var(--red) 7%,var(--surface))}
.status-pill.warn{border-color:color-mix(in srgb,var(--orange) 34%,var(--line));background:color-mix(in srgb,var(--orange) 7%,var(--surface))}
.status-pill.ok{border-color:color-mix(in srgb,var(--green) 28%,var(--line));background:color-mix(in srgb,var(--green) 5%,var(--surface))}

table{width:100%;border-collapse:collapse}
th,td{padding:10px 12px;border-bottom:1px solid var(--line);text-align:left;white-space:nowrap}
th{font-size:11px;color:var(--muted);font-weight:720;background:var(--surface);position:sticky;top:0;z-index:1;text-transform:none}
tbody tr:nth-child(even) td{background:var(--row)}
tr:hover td{background:var(--hover)}
td.num{text-align:right;font-variant-numeric:tabular-nums}
.bad{color:var(--red);font-weight:740}
.muted{color:var(--muted)}
.mini{font-size:11px;color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}

@media(max-width:1180px){
  main{--host-action-reserve:170px}
  .command-grid,.layout{grid-template-columns:1fr}
  .cards{grid-template-columns:repeat(3,1fr)}
  .account-toolbar{grid-template-columns:minmax(220px,1fr) 128px max-content 86px 62px 54px 62px}
  .account-toolbar .spacer{display:none}
  .account-summary-grid{grid-template-columns:repeat(2,minmax(0,1fr))}
  .invalid-auth-row{grid-template-columns:24px minmax(180px,1fr) 78px 60px}
  .invalid-auth-meta,.invalid-auth-reason{grid-column:2 / -1}
}

@media(max-width:720px){
  main{padding:14px 14px 30px;--host-action-reserve:0px}
  .hero{display:block;padding:0 0 12px}
  .controls{justify-content:flex-start;margin-top:10px}
  .cards{grid-template-columns:repeat(2,minmax(0,1fr))}
  .value{font-size:17px}
  .form-row{grid-template-columns:1fr}
  .account-toolbar,.autoban-toolbar{grid-template-columns:1fr 1fr;align-items:center}
  .account-toolbar>input{grid-column:1 / -1;min-width:0;width:100%}
  .autoban-scope,.column-controls{grid-column:1 / -1}
  .column-controls{overflow:auto;padding-bottom:1px}
  .account-summary-grid{grid-template-columns:1fr}
  .account-table-wrap{max-height:540px}
  .autoban-table-wrap{max-height:300px}
  .invalid-auth-panel{width:min(100%,calc(100vw - 14px))}
  .invalid-auth-list{max-height:56vh}
  .invalid-auth-row{grid-template-columns:24px minmax(0,1fr) 1fr 1fr}
  .invalid-auth-row button{width:100%}
  .invalid-auth-main{grid-column:2 / -1}
  .invalid-auth-meta,.invalid-auth-reason{grid-column:2 / -1}
  .insight{grid-template-columns:1fr}
  .insight span:first-child{grid-row:auto}
  .quota-compact{grid-template-columns:28px 1fr 92px;min-width:224px}
  th,td{padding:9px 10px}
}

.log-export-panel{width:min(760px,100%)}
.log-export-grid{grid-template-columns:1fr 1fr}
.log-export-grid .modal-status{grid-column:1 / -1}
.key-summary-table-wrap{max-height:260px}
`
