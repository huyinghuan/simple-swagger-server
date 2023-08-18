
// the following lines will be replaced by docker/configurator, when it runs in a docker-container
window.onload = () => {
  fetch("/api/docs", { method: "POST" }).then((resp) => {
    return resp.json();
  }).then((data) => {
    console.log(data)
    var urls = []
    Object.keys(data).forEach((key) => {
      urls.push({
        url: "/docs/" + data[key],
        name: key
      })
    })
    window.ui = SwaggerUIBundle({
      urls: urls,
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
      filter: true,
    });

    //</editor-fold>
  });
}
