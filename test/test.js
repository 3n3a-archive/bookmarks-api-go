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
        return response;
      })
      .catch(function (error) {
        console.log(error);
      });
}

function allBookmarks(cb) {
      var config = {
        method: 'get',
        url: 'http://localhost:8081/bms',
        headers: { 
          'Content-Type': 'application/json'
        }
      };
      
      axios(config)
      .then(function (response) {
        cb(response.data)
      })
      .catch(function (error) {
        console.log(error);
      });
}

function delBookmark(id) {
    var data = JSON.stringify({
        "ID": id
      });
      
      var config = {
        method: 'delete',
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

function iterDel(data) {
  for (const bm of data) {
    delBookmark(bm.ID)
  }
}

function add_test_bms() {
  // Add all test bookmarks
  for (const item of test_data) {
    addBookmark(item.title, item.description, item.url, item.category)
  }  
}

function del_test_bms() {
  // Delete All test bookmarks
  allBookmarks(iterDel)
}

module.exports = {
  add_test_bms,
  del_test_bms
}