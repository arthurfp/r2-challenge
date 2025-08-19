package main

const swaggerHTML = `<!doctype html><html><head><meta charset="utf-8"/><title>R2 Challenge API</title>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css"></head>
<body><div id="header" style="padding:8px;border-bottom:1px solid #eee;">
  <input id="token" placeholder="Bearer <token>" style="width:60%;padding:6px;" />
  <button onclick="setToken()" style="padding:6px 12px;">Set Token</button>
</div>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
<script>
  let authToken = localStorage.getItem('authToken') || '';
  function setToken(){
    let t = document.getElementById('token').value.trim();
    if(t && !/^Bearer\s+/i.test(t)) { t = 'Bearer ' + t; }
    authToken = t;
    localStorage.setItem('authToken', authToken);
  }
  window.addEventListener('DOMContentLoaded', function(){
    if(authToken){ document.getElementById('token').value = authToken; }
  });
  window.ui = SwaggerUIBundle({
    url:'/swagger.yaml', dom_id:'#swagger-ui',
    requestInterceptor: (req) => { if(authToken){ req.headers['Authorization'] = authToken; } return req; },
    responseInterceptor: (res) => {
      try {
        if(res && res.url && res.status === 200 && /\/v1\/auth\/login$/.test(res.url)){
          const data = typeof res.data === 'string' ? JSON.parse(res.data) : res.data;
          if(data && data.access_token){
            authToken = 'Bearer ' + data.access_token;
            localStorage.setItem('authToken', authToken);
            const el = document.getElementById('token'); if(el){ el.value = authToken; }
          }
        }
      } catch(e){}
      return res;
    }
  });
</script>
</body></html>`
