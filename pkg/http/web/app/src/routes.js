import Vue from 'vue';
import VueRouter from 'vue-router';
import Auth from '@okta/okta-vue'

import Home from './components/Home';
import Login from './components/Login';
import GitHubRepoDetails from './components/GithubRepoDetails';

const REDIRECT_URI = process.env.NODE_ENV == 'production' ? 'https://vue-js-golang-front-end.herokuapp.com' : 'http://localhost:8080'
console.info(REDIRECT_URI)

Vue.use(VueRouter);
Vue.use(Auth, {
  issuer: 'https://dev-509836.oktapreview.com/oauth2/default',
  client_id: '0oagcbm1o6GTTB9Da0h7',
  redirect_uri: REDIRECT_URI + '/implicit/callback',
  scope: 'openid profile email'
})

export default new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', component: Login },
    { path: '/me', component: Home, meta: { requiresAuth: true }},
    { name: 'repo-details', path: '/repo/:id', component: GitHubRepoDetails, meta: { requiresAuth: true } },
    { path: '/implicit/callback', component: Auth.handleCallback() }
  ]
});
