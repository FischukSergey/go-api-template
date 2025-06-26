window.onload = function () {
  // Автоматически определяем текущий сервер
  const currentOrigin = window.location.origin;

  const ui = SwaggerUIBundle({
    url: '/docs/api.yaml',
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout",
    // Устанавливаем текущий сервер как основной
    servers: [
      {
        url: currentOrigin,
        description: "Current server (" + currentOrigin + ")"
      },
      {
        url: "http://localhost:38080",
        description: "Local development server"
      },
      {
        url: "http://localhost:37545",
        description: "Docker development server"
      },
      {
        url: "https://217.12.37.51:37545",  // Замените на ваш VPS домен
        description: "VPS development server"
      },
      {
        url: "https://alex-fisher-team.ru/api/v1",
        description: "Production server"
      }
    ],
    // Используем встроенный oauth2-redirect.html из Swagger UI
    oauth2RedirectUrl: currentOrigin + "/docs/oauth2-redirect.html"
  });

  // OAuth конфигурация
  ui.initOAuth({
    clientId: "swagger-ui",
    realm: "Woman",
    appName: "Woman App API Docs",
    scopes: "openid profile email"
  });
}; 