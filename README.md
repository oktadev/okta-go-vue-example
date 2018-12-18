#### Live Demo: https://vue-js-golang-front-end.herokuapp.com/

---

# Build a Single Page Application with Vue.js and Golang

Single-Page Applications (SPA) has improved the user experience offering rich UI interactions, fast feedback and the feeling that you no longer need to download and install applications on your machine. Browsers are now operating systems and websites are seen as apps. Can you imagine sites like Facebook, Google Maps, Gmail, and so many others refreshing the page on every single interaction? Even though for some web applications it might make sense to not invest time building a SPA, others in which user interactions happen constantly a rich UI could be seen as a determinant factor to the business success.

While in the user’s perspective, SPAs, everything feels like rainbows and unicorns, in the software engineer’s perspective the reality often is the opposite. Well known problems already solved in the back-end like: Authentication, Routing, State management, Data Binding, and so on are now front-end challenges that likely will consume most of your time. Luckily for us many Javascript frameworks were created aiming to help software engineers to craft powerful applications focusing more on the business requirements rather than spending time reinventing the wheel. Frameworks like Vue.js, React, Angular, Ember.js, and so many others provide abstractions and conventions to make our lives easier and the software development experience a enjoyable journey.

## Vue.js

There is no better words to describe Vue.js than the ones from its creator:

> Vue (pronounced `/vjuː/`, like view) is a progressive framework for building user interfaces. It is designed from the ground up to be incrementally adoptable, and can easily scale between a library and a framework depending on different use cases. It consists of an approachable core library that focuses on the view layer only, and an ecosystem of supporting libraries that helps you tackle complexity in large Single-Page Applications.

In a nutshell, the benefits of Vue.js in my opinion are:

Gentle learning curve
`vue-cli` bootstraps your app saving you from the hassle of setting up webpack
`vue-router` and `vuex` are maintained by the `vue.js` team, which means these projects tend to work really really well together as well as tend to evolve together
Community has been growing fast, Vue.js has more stars on Github than javascript frameworks like React and Angular.js
It’s flexible, and by flexible I mean it can be adopted in calm pace, component by component.

## TALK IS CHEAP, SHOW ME THE <CODE /> 
<small>Linus Torvald</small>

In this tutorial we are going to create Single-Page Application to show some love to open source projects hosted on Github. For the front-end, you guessed it right, we are going to use Vue.js and the tooling around it: `vuex`, `vue-cli`, `vuetify`, and `vue-router`. While in the back-end world we are going to use Go to write a REST API and save our data on MongoDB.

### Requirements 

Authentication - A user should be able to identify himself/herself via [Okta’s OpenID Connect (OIDC)](https://developer.okta.com/docs/api/resources/oidc)
When *NOT* Authenticated, a user should be redirected to the Okta’s authentication page
When Authenticated
A user should be able to search for her/his favorite open source projects on Github
A user should be able to favorite the projects returned for the Github search
A user should be able to add notes on any favorited project.

We will use JWT-based authentication when making requests from the our Single-Page Application and [Okta’s JWT Verifier](https://github.com/okta/okta-jwt-verifier-golang) as a middleware on our backend to validate the user’s token for every request.
### Directory structure:

For the sake of simplicity, let’s create the REST API and the SPA in the same project.  Let’s start by creating the project directory into the Go workspace.


```sh
mkdir -p $(go env GOPATH)/src/github.com/{YOUR_GITHUB_USERNAME}/kudo-oos
```

Inside the newly created directory, we will have a structure like this:

```
├── cmd
│   └── db
└── pkg
    ├── core
    ├── http
    │   └── web
    │       └── app
    │           ├── dist
    │           │   ├── css
    │           │   └── js
    │           └── src
    │               ├── assets
    │               ├── components
    │               └── plugins
    ├── kudo
    └── storage
```

No worry about creating all those directories now, we’re going to see each one of them further on in this article. 

To get our SPA off the ground quickly let’s leverage the scaffolding functionality from [vue-cli](https://cli.vuejs.org/). The CLI will prompt a serie options in which the only thing we need to do is to pick the piece of technology we want for our project.

Firstly, install the `vue-cli` by running:

```sh
yarn global add @vue/cli
``` 

Then, create a new vue project:

```sh
mkdir -p pkg/http/web
cd pkg/http/web
vue create app
```

You will be prompted with a serie of questions about the project build details, for this tutorial pick all the default choices. DONE! Congratulations, you have created your Vue.js SPA. Try it by running:

```sh
cd app
yarn install
yarn serve
```

Open this URL: http://localhost:8080 on your browser and you should see the something like this:

![vuetify](https://i.imgur.com/2RkVtV8.png)

Next, let’s see how to make our SPA look modern and responsive using vuetify.

### Meet [Vuetify](https://vuetifyjs.com/en/) 

Vuetify will help us to create good looking SPAs with [Material Design](https://material.io/design/) like the page the user will be redirected after login:


```sh
vue add vuetify
```

Again, you will be prompted with a series of questions, for the sake of simplicity just go with the default choices. Spin up your SPA again to see vuetify in action.

```sh
yarn serve
``` 

![vuetify](https://i.imgur.com/NF2n5I3.png)

## Add Authentication with Okta

Let’s get started, by creating an OIDC application in Okta. Sign up for a forever-free developer account (or log in if you already have one).

![](https://developer.okta.com/assets/blog/vue-crud-node/okta-developer-sign-up-8076b22d5d523a70c1c8a0cef34993103f3a2d01d821d2cc31a1f3ba9798cb08.png)

Once logged in, create a new application by clicking “Add Application”.

![](https://developer.okta.com/assets/blog/vue-crud-node/add-application-024fdfa427033e2e345b48167d8fdef2592f8dcaa464be89487e257a629d39ad.png)

Select the “Single-Page App” platform option.

![](https://developer.okta.com/assets/blog/vue-crud-node/new-application-options-96b1b0faa43717d47faaf99696fe2155e793c5f482b7cbae21747bdf3fd72ba4.png)

The default application settings should be the same as those pictured.

![](https://developer.okta.com/assets/blog/vue-crud-node/okta-application-settings-99892ff83b5e572a4c2b64d3d9b2edde2266d7f7434c18b8e60f6f005d7718e0.png)


Let’s install the Okta Vue SDK, run the following command:

```sh
yarn add @okta/okta-vue
```
### Routes

Create `pkg/http/web/app/src/routes.js` and add the routes:

```js
import Vue from 'vue';
import VueRouter from 'vue-router';
import Auth from '@okta/okta-vue'

import Home from './components/Home';
import Login from './components/Login';
import GitHubRepoDetails from './components/GithubRepoDetails';

Vue.use(VueRouter);
Vue.use(Auth, {
  issuer: {ADD_YOUR_DOMAIN},
  client_id: {ADD_YOUR_CLIENT_ID},
  redirect_uri: 'http://localhost:8080/implicit/callback',
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
```

Make sure to add your `domain` and `client_id` where indicated they can be found on your application overview page in the Okta Developer Console. Calling `Vue.use(Auth, ...)` will inject an `authClient` object into your Vue instance which can be accessed by calling `this.$auth` anywhere inside your Vue instance, which we will use to make sure an user is logged in and/or to force the user to identify himself/herself.

In order to see our Authentication flow working, we will need to create the following files:

```
├── apiClient.js
├── components
│   ├── Footer.vue
│   ├── GithubRepo.vue
│   ├── GithubRepoDetails.vue
│   ├── Home.vue
│   ├── Login.vue
│   └── SearchBar.vue
├── githubClient.js
├── routes.js
└── store.js
```

### API Client

In the `./kudo-oos/pkg/http/web/app/src/apiClient.js` let’s create all methods we need to send requests to our REST API.

```js
import Vue from 'vue';
import axios from 'axios';

const client = axios.create({
  baseURL: 'http://localhost:4444',
  json: true
});

const APIClient =  {
  createKudo(repo) {
    return this.perform('post', '/kudos', repo);
  },

  deleteKudo(repo) {
    return this.perform('delete', `/kudos/${repo.id}`);
  },

  updateKudo(repo) {
    return this.perform('put', `/kudos/${repo.id}`, repo);
  },

  getKudos() {
    return this.perform('get', '/kudos');
  },

  getKudo(repo) {
    return this.perform('get', `/kudo/${repo.id}`);
  },

  async perform (method, resource, data) {
    let accessToken = await Vue.prototype.$auth.getAccessToken()
    return client({
      method,
      url: resource,
      data,
      headers: {
        Authorization: `Bearer ${accessToken}`
      }
    }).then(req => {
      return req.data
    })
  }
}

export default APIClient;
```

Notice that for every single request we inject the user's access token provided by `Vue.prototype.$auth.getAccessToken()` as the `Authorization` header. Futher on, in the back-end, we're going to use this token to make sure the request is valid.

### Vue Components

`./kudo-oos/pkg/http/web/app/src/components/Footer.vue` holds the footer content of our SPA. 

```html
<template>
 <v-footer class="pa-3 white--text" color="teal" absolute>
   <div>
     Developed with ❤️ by {{YOUR_NAME}} &copy; {{ new Date().getFullYear() }}
   </div>
 </v-footer> 
</template>
```

`./kudo-oos/pkg/http/web/app/src/components/GithubRepo.vue` is the component that displays the Github open source project.

```html
<template>
  <v-card >
    <v-card-title primary-title>
      <div class="repo-card-content">
        <h3 class="headline mb-0">
          <router-link :to="{ name: 'repo-details', params: { id: repo.id }}" >{{repo.full_name}}</router-link>
        </h3>
        <div>{{repo.description}}</div>
      </div>
    </v-card-title>
    <v-card-actions>
      <v-chip>
        {{repo.language}}
      </v-chip>
      <v-spacer></v-spacer>
      <v-btn @click.prevent="toggleKudo(repo)"  flat icon color="pink">
        <v-icon v-if="isKudo(repo)">favorite</v-icon>
        <v-icon v-else>favorite_border</v-icon>
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  data() {
    return {}
  },
  props: ['repo'],
  methods: {
    isKudo(repo) {
      return this.$store.getters.isKudo(repo);
    },
    ...mapActions(['toggleKudo'])
  }
}
</script>

<style>
 .repo-card-content {
   height: 90px;
   overflow: scroll;
 }
</style>
```

`./kudo-oos/pkg/http/web/app/src/components/GithubRepoDetails.vue` shows complementary details of the github open source project and allows the user to enter notes, to show some love, about the OSS project. Will be rendered under the `/repo/:id` route.

```html
<template>
  <v-container grid-list-md fluid class="grey lighten-4" >
    <v-layout align-center justify-space-around wrap>
      <v-flex md6>
        <!-- <v-img
          :src="repo.owner.avatar_url"
          :alt="repo.owner.login"
          class="grey darken-4"
          width="200"
        ></v-img> -->
        <h1 class="primary--text">
          <a :href="repo.html_url">{{repo.full_name}}</a>
        </h1>

        <v-chip class="text-xs-center">
          <v-avatar class="teal">
            <v-icon class="white--text">star</v-icon>
          </v-avatar>
          Stars: {{repo.stargazers_count}}
        </v-chip>

        <v-chip class="text-xs-center">
          <v-avatar class="teal white--text">L</v-avatar>
          Language: {{repo.language}}
        </v-chip>

        <v-chip class="text-xs-center">
          <v-avatar class="teal white--text">O</v-avatar>
          Open Issues: {{repo.open_issues_count}}
        </v-chip>

        <v-textarea
          name="input-7-1"
          label="Show some love"
          value=""
          v-model="repo.notes"
          hint="Describe why you love this project"
        ></v-textarea>
        <v-btn @click.prevent="updateKudo(repo)"> Kudo </v-btn>
        <router-link tag="a" to="/me">Back</router-link>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  data() {
    return {
      repo: {}
    }
  },
  watch: {
    '$route': 'fetchData'
  },
  computed: mapGetters(['kudos']),
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      fetch('https://api.github.com/repositories/' + this.$route.params.id)
        .then(response => response.json())
        .then((response) => {
          this.repo = Object.assign(response, this.kudos[this.$route.params.id])
        })
    },
    ...mapActions(['updateKudo'])
  }
}
</script>
```

The `./kudo-oos/pkg/http/web/app/src/components/Home.vue` is the father of all components.

```html
<template>
 <div>  
   <SearchBar defaultQuery='okta' v-on:search-submitted="githubQuery" />
   <v-container grid-list-md fluid class="grey lighten-4" >
        <v-tabs
       slot="extension"
       v-model="tabs"
       centered
       color="teal"
       text-color="white"
       slider-color="white"
     >
       <v-tab class="white--text" :key="2">
         KUDOS
       </v-tab>
       <v-tab class="white--text" :key="1">
         SEARCH
       </v-tab>
     </v-tabs>
       <v-tabs-items style="width:100%" v-model="tabs">
         <v-tab-item :key="2">
           <v-layout row wrap>
             <v-flex v-for="kudo in allKudos" :key="kudo.id" md4 >
               <GitHubRepo :repo="kudo" />
             </v-flex>
           </v-layout>
         </v-tab-item>
         <v-tab-item :key="1">
           <v-layout row wrap>
             <v-flex v-for="repo in repos" :key="repo.id" md4>
               <GitHubRepo :repo="repo" />
             </v-flex>
           </v-layout>
         </v-tab-item>
       </v-tabs-items>
   </v-container>
 </div>
</template>

<script>
import SearchBar from './SearchBar.vue'
import GitHubRepo from './GithubRepo.vue'
import githubClient from '../githubClient'
import { mapMutations, mapGetters, mapActions } from 'vuex'

export default {
 name: 'Home',
 components: { SearchBar, GitHubRepo },
 data() {
   return {
     tabs: 0
   }
 },
 computed: mapGetters(['allKudos', 'repos']),
 created() {
   this.getKudos();
 },
 methods: {
   githubQuery(query) {
     this.tabs = 1;
     githubClient
       .getJSONRepos(query)
       .then(response => this.resetRepos(response.items) )
   },
   ...mapMutations(['resetRepos']),
   ...mapActions(['getKudos']),
 },
}
</script>

<style>
.v-tabs__content {
  padding-bottom: 2px;
}
</style>
```

`./kudo-oos/pkg/http/web/app/src/components/Login.vue` it’s the simplest component we have, it only have a login button.

```html
<template>
 <v-app id="inspire">
   <v-content>
     <v-container fluid fill-height>
       <v-layout align-center justify-center>
         <v-flex xs12 sm8 md4>
           <v-card class="elevation-12">
             <v-toolbar dark color="teal">
               <v-toolbar-title justify-center>Login</v-toolbar-title>
             </v-toolbar>
             <v-card-text>
               <v-btn @click.prevent="login" color="primary">Sign in with Okta</v-btn>
             </v-card-text>
           </v-card>
         </v-flex>
       </v-layout>
     </v-container>
   </v-content>
 </v-app>
</template>

<script>
export default {
 data() {
   return {};
 },
 async mounted() {
   const isAuthenticated = await this.$auth.isAuthenticated();
   isAuthenticated && this.$router.push('/me');
 },
 methods: {
   login () {
     this.$auth.loginRedirect('/me')
   }
 }
}
</script>
```

`./kudo-oos/pkg/http/web/app/src/components/SearchBar.vue` our last component.

```html
<template>
   <v-toolbar dark color="teal">
     <v-spacer></v-spacer>
     <v-text-field
       solo-inverted
       flat
       hide-details
       label="Search for your OOS project on Github + Press Enter"
       prepend-inner-icon="search"
       v-model="query"
       @keyup.enter="onSearchSubmition"
     ></v-text-field>
     <v-spacer></v-spacer>
     <button @click.prevent="logout">Logout</button>
   </v-toolbar>
</template>

<script>
export default {
 data() {
   return {
     query: null,
   };
 },
 props: ['defaultQuery'],
 methods: {
   onSearchSubmition() {
     this.$emit('search-submitted', this.query);
   },
   async logout () {
     await this.$auth.logout()
     this.$router.push('/')
   }
 }
}
</script>
```
### Github Client

The search feature is backed by the Github API. Lets create `./kudo-oos/pkg/http/web/app/src/githubClient.js` and place the methods we need there.

```js
const API_URL = "https://api.github.com/search/repositories"
export default {
 getJSONRepos(query) {
   return fetch(`${API_URL}?q=` + query).then(response => response.json());
 }
}
```
### [Vuex](https://vuex.vuejs.org/)

> Vuex is a state management pattern + library for Vue.js applications. It serves as a centralized store for all the components in an application, with rules ensuring that the state can only be mutated in a predictable fashion. It also integrates with Vue's official devtools extension to provide advanced features such as zero-config time-travel debugging and state snapshot export / import.

Even though our SPA isn’t a large web application, in this project we are using `vuex` to manage the application state here how our `./kudo-oos/pkg/http/web/app/src/store.js` looks like: 

```js
import Vue from 'vue';
import Vuex from 'vuex';

import APIClient from './apiClient';

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    kudos: {},
    repos: [],
  },
  mutations: {
    resetRepos (state, repos) {
      state.repos = repos;
    },
    resetKudos(state, kudos) {
      state.kudos = kudos;
    }
  },
  getters: {
    allKudos(state) {
      return Object.values(state.kudos);
    },
    kudos(state) {
      return state.kudos;
    },
    repos(state) {
      return state.repos;
    },
    isKudo(state) {
      return (repo)=> {
        return !!state.kudos[repo.id];
      };
    }
  },
  actions: {
    getKudos ({commit}) {
      APIClient.getKudos().then((data) => {
        commit('resetKudos', data.reduce((acc, kudo) => { 
                               return {[kudo.id]: kudo, ...acc}
                             }, {}))
      })
    },
    updateKudo({ commit, state }, repo) {
      const kudos = { ...state.kudos, [repo.id]: repo };
      
      return APIClient
        .updateKudo(repo)
        .then(() => {
          commit('resetKudos', kudos)
        });
    },
    toggleKudo({ commit, state }, repo) {
      if (!state.kudos[repo.id]) {
        return APIClient
          .createKudo(repo)
          .then(kudo => commit('resetKudos', { [kudo.id]: kudo, ...state.kudos }))
      }

      const kudos = Object.entries(state.kudos).reduce((acc, [repoId, kudo]) => {
                      return (repoId == repo.id) ? acc
                                                 : { [repoId]: kudo, ...acc };
                    }, {});

      return APIClient
        .deleteKudo(repo)
        .then(() => commit('resetKudos', kudos));
    }
  }
});

export default store;
```
### Entrypoint 

Ok, our router, store and components are in place let's modify `./kudo-oos/pkg/http/web/app/src/main.js` to properly initiate our SPA:

```js
import '@babel/polyfill'
import Vue from 'vue'
import './plugins/vuetify'
import App from './App.vue'
import store from './store'
import router from './routes'

Vue.config.productionTip = process.env.NODE_ENV == 'production';

router.beforeEach(Vue.prototype.$auth.authRedirectGuard())

new Vue({
 store,
 router,
 render: h => h(App)
}).$mount('#app')
```

Note that we are calling `router.beforeEach(Vue.prototype.$auth.authRedirectGuard())` which will look for routes tagged with  `meta: {requiresAuth: true}` and redirects the user to the authentication flow if they are not authenticated.

Here’s the page the user should see when not authenticated:
Then once he/she clicks in the login button, he/she should be redirected to the Okta’s login page:
And after a successful login, the user is redirected back to our application:

## REST API

Now that users can securely authenticate, you can build the REST API. 

### Directory Structure

Here’s how the directory structure looks like:

```sh
tree -I "vendor|web"

├── Gopkg.lock
├── Gopkg.toml
├── Makefile
├── Procfile
├── cmd
│   ├── db
│   │   └── setup.go
│   └── main.go
├── docker-compose.yml
└── pkg
    ├── core
    │   ├── kudo.go
    │   └── repository.go
    ├── http
    │   ├── handlers.go
    │   └── middlewares.go
    ├── kudo
    │   └── service.go
    └── storage
        ├── mongo.go
        ├── mongo_test.go
        └── storage_suite_test.go
```

### Dependencies

Let’s start by downloading the essential dependencies:

```sh
dep init
dep ensure -add github.com/okta/okta-jwt-verifier-golang
dep ensure -add github.com/rs/cors
dep ensure -add github.com/globalsign/mgo
```

### Add MongoDB Repository Implementation

Go ahead and create an interface to represent a repository of data.

`./kudo-oos/pkg/core/kudo.go` defines the struct which represents a open source project. 

```go
package core

// Kudo represents a oos kudo.
type Kudo struct {
  UserID      string `json:"user_id" bson:"userId"`
  RepoID      string `json:"id" bson:"repoId"`
  RepoName    string `json:"full_name" bson:"repoName"`
  RepoURL     string `json:"html_url" bson:"repoUrl"`
  Language    string `json:"language" bson:"language"`
  Description string `json:"description" bson:"description"`
  Notes       string `json:"notes" bson:"notes"`
}
```

Where as `./kudo-oos/pkg/core/repository.go` implements the interface which represents our repository API.

```go
package core
// Repository defines the API a repository implementation should follow.
type Repository interface {
  Find(id string) (*Kudo, error)
  FindAll(selector map[string]interface{}) ([]*Kudo, error)
  Delete(kudo *Kudo) error
  Update(kudo *Kudo) error
  Create(kudo ...*Kudo) error
  Count() (int, error)
}
```

Now, let’s add a MongoDB implementation of the repository: `./kudo-oos/pkg/storage/mongo.go`

```go
package storage

import (
  "log"
  "os"

  "github.com/globalsign/mgo"
  "github.com/globalsign/mgo/bson"
  "github.com/{YOUR_GITHUB_USERNAME}/kudo-oos/pkg/core"
)

const (
  collectionName = "kudos"
)

func GetCollectionName() string {
  return collectionName
}

type MongoRepository struct {
  logger  *log.Logger
  session *mgo.Session
}

// Find fetches a kudo from mongo according to the query criteria provided.
func (r MongoRepository) Find(repoID string) (*core.Kudo, error) {
  session := r.session.Copy()
  defer session.Close()
  coll := session.DB("").C(collectionName)

  var kudo core.Kudo
  err := coll.Find(bson.M{"repoId": repoID, "userId": kudo.UserID}).One(&kudo)
  if err != nil {
    r.logger.Printf("error: %v\n", err)
    return nil, err
  }
  return &kudo, nil
}

// FindAll fetches kudos from the database.
func (r MongoRepository) FindAll(selector map[string]interface{}) ([]*core.Kudo, error) {
  session := r.session.Copy()
  defer session.Close()
  coll := session.DB("").C(collectionName)

  var kudos []*core.Kudo
  err := coll.Find(selector).All(&kudos)
  if err != nil {
    r.logger.Printf("error: %v\n", err)
    return nil, err
  }
  return kudos, nil
}

// Delete deletes a kudo from mongo according to the query criteria provided.
func (r MongoRepository) Delete(kudo *core.Kudo) error {
  session := r.session.Copy()
  defer session.Close()
  coll := session.DB("").C(collectionName)

  return coll.Remove(bson.M{"repoId": kudo.RepoID, "userId": kudo.UserID})
}

// Update updates an kudo.
func (r MongoRepository) Update(kudo *core.Kudo) error {
  session := r.session.Copy()
  defer session.Close()
  coll := session.DB("").C(collectionName)

  return coll.Update(bson.M{"repoId": kudo.RepoID, "userId": kudo.UserID}, kudo)
}

// Create kudos in the database.
func (r MongoRepository) Create(kudos ...*core.Kudo) error {
  session := r.session.Copy()
  defer session.Close()
  coll := session.DB("").C(collectionName)

  for _, kudo := range kudos {
    _, err := coll.Upsert(bson.M{"repoId": kudo.RepoID, "userId": kudo.UserID}, kudo)
    if err != nil {
      return err
    }
  }

  return nil
}

// Count counts documents for a given collection
func (r MongoRepository) Count() (int, error) {
  session := r.session.Copy()
  defer session.Close()
  coll := session.DB("").C(collectionName)
  return coll.Count()
}

// NewMongoSession dials mongodb and creates a session.
func newMongoSession() (*mgo.Session, error) {
  mongoURL := os.Getenv("MONGO_URL")
  if mongoURL == "" {
    log.Fatal("MONGO_URL not provided")
  }
  return mgo.Dial(mongoURL)
}

func newMongoRepositoryLogger() *log.Logger {
  return log.New(os.Stdout, "[mongoDB] ", 0)
}

func NewMongoRepository() core.Repository {
  logger := newMongoRepositoryLogger()
  session, err := newMongoSession()
  if err != nil {
    logger.Fatalf("Could not connect to the database: %v\n", err)
  }

  return MongoRepository{
    session: session,
    logger:  logger,
  }
}
```

### Kudo

Before we create our handlers, let's create a piece of code that knows how to handle incoming requests payload as well as interpret them to perform CRUD operations against MongoDB.

`./kudo-oos/pkg/kudo/service.go`

```go
package kudo

import (
  "strconv"

  "github.com/{YOUR_GITHUB_USERNAME}/kudo-oos/pkg/core"
)

type GitHubRepo struct {
  RepoID      int64  `json:"id"`
  RepoURL     string `json:"html_url"`
  RepoName    string `json:"full_name"`
  Language    string `json:"language"`
  Description string `json:"description"`
  Notes       string `json:"notes"`
}

type Service struct {
  userId string
  repo   core.Repository
}

func (s Service) GetKudos() ([]*core.Kudo, error) {
  return s.repo.FindAll(map[string]interface{}{"userId": s.userId})
}

func (s Service) CreateKudoFor(githubRepo GitHubRepo) (*core.Kudo, error) {
  kudo := s.githubRepoToKudo(githubRepo)
  err := s.repo.Create(kudo)
  if err != nil {
    return nil, err
  }
  return kudo, nil
}

func (s Service) UpdateKudoWith(githubRepo GitHubRepo) (*core.Kudo, error) {
  kudo := s.githubRepoToKudo(githubRepo)
  err := s.repo.Create(kudo)
  if err != nil {
    return nil, err
  }
  return kudo, nil
}

func (s Service) RemoveKudo(githubRepo GitHubRepo) (*core.Kudo, error) {
  kudo := s.githubRepoToKudo(githubRepo)
  err := s.repo.Delete(kudo)
  if err != nil {
    return nil, err
  }
  return kudo, nil
}

func (s Service) githubRepoToKudo(githubRepo GitHubRepo) *core.Kudo {
  return &core.Kudo{
    UserID:      s.userId,
    RepoID:      strconv.Itoa(int(githubRepo.RepoID)),
    RepoName:    githubRepo.RepoName,
    RepoURL:     githubRepo.RepoURL,
    Language:    githubRepo.Language,
    Description: githubRepo.Description,
    Notes:       githubRepo.Notes,
  }
}

func NewService(repo core.Repository, userId string) Service {
  return Service{
    repo:   repo,
    userId: userId,
  }
}
```

### HTTP Handlers

Our REST API needs to expose the following endpoints, .

```
 # Fetches all open source projects favorited by the user
GET /kudos
# Fetches a favorited open source project by id
GET /kudos/:id 
# Creates (or favorites)  a open source project for the logged in user
POST /kudos
# Updates  a favorited open source project
PUT /kudos/:id
# Deletes (or unfavorites) a favorited open source project
DELETE /kudos/:id
```

Let’s create `./kudo-oos/pkg/http/handlers.go`

```go
package http

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "strconv"

  "github.com/julienschmidt/httprouter"
  "github.com/{YOUR_GITHUB_USERNAME}/kudo-oos/pkg/core"
  "github.com/{YOUR_GITHUB_USERNAME}/kudo-oos/pkg/kudo"
)

type Service struct {
  repo   core.Repository
  Router http.Handler
}

func New(repo core.Repository) Service {
  service := Service{
    repo: repo,
  }

  router := httprouter.New()
  router.GET("/kudos", service.Index)
  router.POST("/kudos", service.Create)
  router.DELETE("/kudos/:id", service.Delete)
  router.PUT("/kudos/:id", service.Update)

  service.Router = UseMiddlewares(router)

  return service
}

func (s Service) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  service := kudo.NewService(s.repo, r.Context().Value("userId").(string))
  kudos, err := service.GetKudos()

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(kudos)
}

func (s Service) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  service := kudo.NewService(s.repo, r.Context().Value("userId").(string))
  payload, _ := ioutil.ReadAll(r.Body)

  githubRepo := kudo.GitHubRepo{}
  json.Unmarshal(payload, &githubRepo)

  kudo, err := service.CreateKudoFor(githubRepo)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(kudo)
}

func (s Service) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  service := kudo.NewService(s.repo, r.Context().Value("userId").(string))

  repoID, _ := strconv.Atoi(params.ByName("id"))
  githubRepo := kudo.GitHubRepo{RepoID: int64(repoID)}

  _, err := service.RemoveKudo(githubRepo)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
}

func (s Service) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  service := kudo.NewService(s.repo, r.Context().Value("userId").(string))
  payload, _ := ioutil.ReadAll(r.Body)

  githubRepo := kudo.GitHubRepo{}
  json.Unmarshal(payload, &githubRepo)

  kudo, err := service.UpdateKudoWith(githubRepo)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(kudo)
}
```

### Verify Your JWT

This is the most crucial component of your REST API server. Without this middleware any user can perform CRUD operations on our database. In case no authorization header is present or when token is invalid, then API call will be aborted an error will be returned to the client.

Create `./kudo-oos/pkg/http/middlewares.go` and paste the following code:

```go
package http

import (
  "context"
  "log"
  "net/http"
  "strings"

  jwtverifier "github.com/okta/okta-jwt-verifier-golang"
  "github.com/rs/cors"
)

func OktaAuth(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    accessToken := r.Header["Authorization"]
    jwt, err := validateAccessToken(accessToken)
    if err != nil {
      w.WriteHeader(http.StatusForbidden)
      w.Write([]byte(err.Error()))
      return
    }
    ctx := context.WithValue(r.Context(), "userId", jwt.Claims["sub"].(string))
    h.ServeHTTP(w, r.WithContext(ctx))
  })
}

func validateAccessToken(accessToken []string) (*jwtverifier.Jwt, error) {
  parts := strings.Split(accessToken[0], " ")
  jwtVerifierSetup := jwtverifier.JwtVerifier{
    Issuer:           "{DOMAIN}",
    ClaimsToValidate: map[string]string{"aud": "api://default", "cid": "{CLIENT_ID}"},
  }
  verifier := jwtVerifierSetup.New()
  return verifier.VerifyIdToken(parts[1])
}

func JSONApi(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    h.ServeHTTP(w, r)
  })
}

func AccsessLog(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Printf("%s: %s", r.Method, r.RequestURI)
    h.ServeHTTP(w, r)
  })
}

func Cors(h http.Handler) http.Handler {
  corsConfig := cors.New(cors.Options{
    AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
    AllowedMethods: []string{"POST", "PUT", "GET", "PATCH", "OPTIONS", "HEAD", "DELETE"},
    Debug:          true,
  })
  return corsConfig.Handler(h)
}

func UseMiddlewares(h http.Handler) http.Handler {
  h = JSONApi(h)
  h = OktaAuth(h)
  h = Cors(h)
  return AccsessLog(h)
}
```
As you can see, the middleware `OktaAuth` uses [okta-jwt-verifier-golang](github.com/okta/okta-jwt-verifier-golang) to validate the user's access token. 

### Entrypoint

`./kudo-oos/pkg/cmd/main.go` will spin up our Go server.

```go
package main

import (
  "log"
  "net/http"
  "os"

  web "github.com/{YOUR_GITHUB_USERNAME}/kudo-oos/pkg/http"
  "github.com/{YOUR_GITHUB_USERNAME}/kudo-oos/pkg/storage"
)

func main() {
  httpPort := os.Getenv("PORT")

  repo := storage.NewMongoRepository()
  webService := web.New(repo)

  log.Printf("Running on port %s\n", httpPort)
  log.Fatal(http.ListenAndServe(httpPort, webService.Router))
}
```

## Running our Applications

### Makefile

Let's create a `Makefile`

```
setup: run_services
	@go run ./cmd/db/setup.go

run_services:
	@docker-compose up --build -d

run_server:
	@MONGO_URL=mongodb://mongo_user:mongo_secret@0.0.0.0:27017/kudos PORT=:4444 go run cmd/main.go

run_client:
    @/bin/bash -c "cd $$GOPATH/src/github.com/klebervirgilio/kudo-oos/pkg/http/web/app && yarn serve"
```

### Dockerfile

As you can see you are using docker to run MongoDB, let's create `docker-compose.yml`

```yaml
version: '3'
services:
  mongo:
    image: mongo
    restart: always
    ports:
     - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: mongo_user
      MONGO_INITDB_ROOT_PASSWORD: mongo_secret
```

### Run

Great! Now all we have to do is:

```sh
make setup
make run_server
make run_client
```

Assuming all went well, you should have Go REST API listening to the port `0.0.0.0:4444` and the SPA serving files on `http://localhost:8080`.

## Conclusions

Vue.js is really powerful and at the same time straightforward framework, its adoption has been growing and the community is becoming stronger. In this tutorial we covered the full cycle of a Single-Page Application development.
To learn more about Vue.js head over to https://vuejs.org or check out these other great resources from the @oktadev team:
The Ultimate Guide to Progressive Web Applications
The Lazy Developer’s Guide to Authentication with Vue.js
Build a Cryptocurrency Comparison Site with Vue.js
Let me know your thoughts  in the comments and feel free to to make any questions, and as always, follow @oktadev on Twitter to see all the cool content our dev team is creating.
