'use strict';(function(e,c,m){if(e.JSON){var G=function(a,d){var b=(e.XMLHttpRequest?function(){return new e.XMLHttpRequest}:function(){return new e.ActiveXObject("Microsoft.XMLHTTP")})();b.open("GET",a,!0);b.onreadystatechange=function(){b&&4==b.readyState&&(200==b.status?d(null,e.JSON.parse(b.responseText)):d(b))};b.send({}.data||"");return b},x=function(a){return String(a).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;").replace(/"/g,"&quot;")},y=function(a,d){var b;return function(){b&&
e.clearTimeout(b);b=e.setTimeout(a,d)}},z=function(a){27==a.keyCode&&c.activeElement&&c.activeElement.blur();var d=c.activeElement.tabIndex;if(30==d)c.querySelector('[tabIndex="6"]').focus(),a.preventDefault();else if(!(6>d)){if(38==a.keyCode){var b=c.querySelector('[tabindex="'+(d-1)+'"]');b.focus();a.preventDefault()}40==a.keyCode&&(b=c.querySelector('[tabindex="'+(d+1)+'"]'),b.focus(),a.preventDefault())}return!a.metaKey&&!a.altKey&&!a.ctrlKey&&46<a.keyCode&&(112>a.keyCode||145<a.keyCode)&&(91>
a.keyCode||95<a.keyCode)},h=function(a,d){var b=[];a.g&&b.push("g="+a.g);a.q&&!1!==d&&1<a.p+d&&b.push("p="+(a.p+d));a.q&&b.push("q="+encodeURIComponent(a.q).replace(/%20/g,"+"));return b.length?"/?"+b.join("&"):"/"},A=function(a){return a.q?a.q+" | Common Search":"Common Search"},B=c.getElementById("f"),g=c.getElementById("q"),H=c.getElementById("hits"),r=c.getElementById("pager"),C=c.getElementById("dbg"),k=c.getElementById("g").childNodes[0],l=c.getElementById("logo"),I=c.getElementsByTagName("title")[0],
t=c.body.className,n=e.JSON.parse(B.getAttribute("data-init")||"{}"),f=n,p=null,u=function(){return{q:g.value.trim(),p:parseInt(r.getAttribute("data-page"),10)||1,g:k.value}},D=function(){"hits"!=t&&(t=c.body.className="hits")},v=function(){};e.history&&e.history.pushState&&(v=function(){var a=f,d=A(a);h(a,0)!=f&&e.history.pushState({s:a},d,h(a,!1));I.innerHTML=x(d);c.title=d},e.history.state||e.history.replaceState({s:f},A(f),h(f,0)));var J=y(v,2E3),E=function(a,d){var b=5,c="";d.c&&(c+="<div id='c'>About "+
d.c+" results</div>");for(var e=0;e<(d.h||[]).length;e++)var f=d.h[e],c=c+("<div class='r'><h3><a href='"+f.u+"' tabindex='"+(b+=1)+"'>"+f.t+"</a></h3><div class='u'><a href='"+f.u+"' tabIndex='-1'>"+f.u.replace(/(.*?:\/\/)(([^\/]+)(\/.+)?).*/,"$2").substring(0,100)+"</a></div><div class='s'>"+f.s+"</div></div>");!c&&a.q&&(c="<div class='z'>We didn't find any results for this search, sorry!</div>");H.innerHTML=c;b="";a.a&&1<a.a&&(b='<a href="'+h(a,-1)+'">&laquo; Previous</a>');d.d&&(b+='<a href="'+
h(a,1)+'">Next &raquo;</a>');r.innerHTML=b;r.setAttribute("data-page",a.a);d.t?(b=d.t,C.innerHTML="Text: <span>"+b.tq+" / "+b.tr+"us</span><br/>Docs: <span>"+b.dq+" / "+b.dr+"us</span><br/>Total: <span>"+b.o+"us</span><br/>"):C.innerHTML=""},K=function(a,d){n=a;a.q?(D(),p&&p.abort(),p=G("/api/search"+h(a,0).substring(1),function(b,c){p=null;b||(c.r&&d?e.location=c.r:E(a,c))})):E(a,{})},q=function(a,c){f=u();a&&!c&&v();h(f,0)!=h(n,0)&&(a||c||"full"==t||J(),f.p=1,K(f,a))},w=y(function(){q(!1,!1)},150);
n.g||(k.value=((m.b?m.b[0]:m.language||(m||{}).userLanguage)||"en").substring(0,2).toLowerCase(),k.value||(k.value="en"));g.onkeydown=function(a){z(a)&&D();w();a.stopPropagation()};B.onsubmit=function(){q(!0,!1);return!1};g.onchange=function(){f.q!=u().q&&w()};g.onmouseup=g.onchange;c.body.onkeydown=function(a){z(a)&&(g.value=g.value.replace(/\s$/,"")+" ",g.focus(),e.scrollTo(0,0))};var F=function(a){0>l.href.indexOf("?g=")?l.href+="?g="+a:l.href=l.href.slice(0,-2)+a};k.onchange=function(){F(k.value);
f.g!=u().g&&w()};e.onpopstate=function(a){a.state&&a.state.s&&(g.value=x(a.state.s.q),q(!0,!0))};F(k.value);q(!1,!1)}})(window,document,navigator);
