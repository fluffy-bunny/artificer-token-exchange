{{define "views/graphiql/index"}}
{{template "html_begin" .}}
{{template "header-graphiql" .}}
<body>
 
  <div id="splash">
    Loading&hellip;
  </div>
  <script src="//cdn.jsdelivr.net/es6-promise/4.0.5/es6-promise.auto.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/react/umd/react.production.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/react-dom/umd/react-dom.production.min.js"></script>
  <script src="//unpkg.com/graphiql/graphiql.min.js"></script>
  <script>
        function getCookieValue(name) {
            const nameString = name + "="
        
            const value = document.cookie.split(";").filter(item => {
            return item.includes(nameString)
            })
        
            if (value.length) {
            return value[0].substring(nameString.length, value[0].length)
            } else {
            return ""
            }
        }
      // Parse the search string to get url parameters.
      var search = window.location.search;
      var parameters = {};
      search.substr(1).split('&').forEach(function (entry) {
        var eq = entry.indexOf('=');
        if (eq >= 0) {
          parameters[decodeURIComponent(entry.slice(0, eq))] =
            decodeURIComponent(entry.slice(eq + 1));
        }
      });

      // if variables was provided, try to format it.
      if (parameters.variables) {
        try {
          parameters.variables =
            JSON.stringify(JSON.parse(parameters.variables), null, 2);
        } catch (e) {
          // Do nothing, we want to display the invalid JSON as a string, rather
          // than present an error.
        }
      }

      // When the query and variables string is edited, update the URL bar so
      // that it can be easily shared
      function onEditQuery(newQuery) {
        parameters.query = newQuery;
        updateURL();
      }
      function onEditVariables(newVariables) {
        parameters.variables = newVariables;
        updateURL();
      }
      function onEditOperationName(newOperationName) {
        parameters.operationName = newOperationName;
        updateURL();
      }
      function updateURL() {
        var newSearch = '?' + Object.keys(parameters).filter(function (key) {
          return Boolean(parameters[key]);
        }).map(function (key) {
          return encodeURIComponent(key) + '=' +
            encodeURIComponent(parameters[key]);
        }).join('&');
        history.replaceState(null, null, newSearch);
      }

       function graphQLFetcher(graphQLParams) {
           let csrf = getCookieValue('_csrf');
          // This example expects a GraphQL server at the path /graphql.
          // Change this to point wherever you host your GraphQL server.
          return fetch('/api/vi/graphql', {
            method: 'post',
            credentials: 'include',
            headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
              'X-Csrf-Token': csrf
            },
            body: JSON.stringify(graphQLParams),
          }).then(function (response) {
            return response.text();
          }).then(function (responseBody) {
            try {
              return JSON.parse(responseBody);
            } catch (error) {
              return responseBody;
            }
          });
        }

      // Render <GraphiQL /> into the body.
      ReactDOM.render(
        React.createElement(GraphiQL, {
          fetcher: graphQLFetcher,
          query: parameters.query,
          variables: parameters.variables,
          operationName: parameters.operationName,
          onEditQuery: onEditQuery,
          onEditVariables: onEditVariables,
          onEditOperationName: onEditOperationName
        }),
        document.body,
      );
  </script>
     
</body>
    
{{template "html_end" .}}
{{end}}