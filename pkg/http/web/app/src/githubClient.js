const API_URL = "https://api.github.com/search/repositories"
export default {
  getJSONRepos(query) {
    return fetch(`${API_URL}?q=` + query).then(response => response.json());
  }
}