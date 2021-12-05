var axios = require('axios');
var test_data = require('./test.json')

function addBookmark(title, description, url, category) {
    var data = JSON.stringify({
        "Title": title,
        "Description": description,
        "URL": url,
        "Category": category
      });
      
      var config = {
        method: 'post',
        url: 'http://localhost:8081/bm',
        headers: { 
          'Content-Type': 'application/json'
        },
        data : data
      };
      
      axios(config)
      .then(function (response) {
        console.log(JSON.stringify(response.data));
      })
      .catch(function (error) {
        console.log(error);
      });
}

for (const item of test_data) {
    addBookmark(item.title, item.description, item.url, item.category)
}