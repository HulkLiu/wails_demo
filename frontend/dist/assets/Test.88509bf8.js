import{_ as g,K as T,L as y,M as C,P as x,Q as S,R as N,T as L,a as _,b as u,n as b,w as c,o as r,e as i,c as M,s as P,m as U,t as B,F as k}from"./index.a5e916c4.js";const w={components:{NLayout:T,NLayoutSider:y,NLayoutContent:C,NMenu:x,NTabs:S,NTabPane:N},setup(){L();const l=_("home"),s=[{label:"Home",key:"home",path:"/home"},{label:"Profile",key:"profile",path:"/about"},{label:"Settings",key:"settings",path:"/videoList"},{label:"Settings1",key:"settings1",path:"/setting"},{label:"Settings2",key:"settings2",path:"/test"}],n=_([{label:"Home",name:"home"},{label:"Profile",name:"profile"},{label:"Settings",name:"settings"}]);return{activeTab:l,menuOptions:s,tabPanes:n,handleMenuSelect:e=>{n.value.find(o=>o.name===e)||n.value.push({label:e.charAt(0).toUpperCase()+e.slice(1),name:e}),l.value=e},handleTabChange:e=>{l.value=e,console.log(e)},handleTabClose:e=>{const o=n.value.findIndex(d=>d.name===e);o!==-1&&(n.value.splice(o,1),l.value===e&&(n.value.length?l.value=n.value[Math.max(0,o-1)].name:l.value=""))}}}},F={style:{padding:"16px"}};function H(l,s,n,a,p,m){const e=u("n-menu"),o=u("n-layout-sider"),d=u("n-tab-pane"),v=u("n-tabs"),h=u("n-layout-content"),f=u("n-layout");return r(),b(f,{"has-sider":""},{default:c(()=>[i(o,{bordered:""},{default:c(()=>[i(e,{options:a.menuOptions,value:a.activeTab,"onUpdate:value":[s[0]||(s[0]=t=>a.activeTab=t),a.handleMenuSelect]},null,8,["options","value","onUpdate:value"])]),_:1}),i(h,null,{default:c(()=>[i(v,{type:"card",value:a.activeTab,"onUpdate:value":[s[1]||(s[1]=t=>a.activeTab=t),a.handleTabChange],closable:"",onClose:a.handleTabClose},{default:c(()=>[(r(!0),M(k,null,P(a.tabPanes,t=>(r(),b(d,{key:t.name,name:t.name,label:t.label},{default:c(()=>[U("div",F,"Content of "+B(t.label),1)]),_:2},1032,["name","label"]))),128))]),_:1},8,["value","onUpdate:value","onClose"])]),_:1})]),_:1})}const R=g(w,[["render",H]]);export{R as default};