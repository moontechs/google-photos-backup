<script setup>
import { useRouter } from 'vue-router'
import { inject } from 'vue'

const router = useRouter()
const db = inject('db')

const logout = () => {
  db.signout()
  router.push({name: 'login'})
}

</script>

<template>
    <div class="page">
    <header class="navbar navbar-expand-sm navbar-light d-print-none">
      <div class="container-xl">
        <h1 class="navbar-brand navbar-brand-autodark d-none-navbar-horizontal pe-0 pe-md-3">
          <a href="#">
            <img src="..." width="110" height="32" alt="Google backup" class="navbar-brand-image" />
          </a>
        </h1>

        <div class="navbar-nav flex-row order-md-last">
          <div class="nav-item">
            <a href="/" class="nav-link d-flex lh-1 text-reset p-0">
              <div class="d-none d-xl-block ps-2">
                <button v-on:click="logout" class="btn btn-primary">Logout</button>
              </div>
            </a>
          </div>
        </div>
      </div>
    </header>
    <div class="page-wrapper">
      <div class="page-body">
        <div class="container-xl">
          <div class="row row-cards">
            <div class="col-lg-4">
              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">New client</h3>
                </div>
                <div class="card-body">
                  <div class="mb-3">
                    <div class="form-label">Client ID</div>
                    <input v-model="newClientId" type="text" class="form-control" autocomplete="off">
                  </div>
                  <div class="mb-3">
                    <div class="form-label">Client secret</div>
                    <input v-model="newClientSecret" type="password" class="form-control">
                  </div>
                  <div class="mt-2">
                    <a v-on:click="addNewClient" href="#" class="btn btn-primary w-100">
                      Add
                    </a>
                  </div>
                </div>
              </div>
            </div>
<!--            <div class="col-lg-4">-->
<!--              <div class="card card-md">-->
<!--&lt;!&ndash;                <div class="card-header">&ndash;&gt;-->
<!--&lt;!&ndash;                  <h3 class="card-title">Add new client</h3>&ndash;&gt;-->
<!--&lt;!&ndash;                </div>&ndash;&gt;-->
<!--                <div class="card-body" style="height: 10rem">-->
<!--                  <div class="mb-3">-->
<!--                    <label class="form-label">Client ID</label>-->
<!--                    <input v-model="clientId" type="text" class="form-control" placeholder="" autocomplete="off">-->
<!--                  </div>-->
<!--                  <div class="mb-2">-->
<!--                    <label class="form-label">-->
<!--                      Client Secret-->
<!--                    </label>-->
<!--                    <div class="input-group input-group-flat">-->
<!--                      <input v-model="clientSecret" type="password" class="form-control"  placeholder="Your password"  autocomplete="off">-->
<!--                      <span class="input-group-text">-->
<!--                  <a href="#" class="link-secondary" title="Show password" data-bs-toggle="tooltip">&lt;!&ndash; Download SVG icon from http://tabler-icons.io/i/eye &ndash;&gt;-->
<!--                    <svg xmlns="http://www.w3.org/2000/svg" class="icon" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" /><path d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6" /></svg>-->
<!--                  </a>-->
<!--                </span>-->
<!--                    </div>-->
<!--                  </div>-->

<!--                  <div class="alert alert-danger" role="alert" v-show="clientErrorMessage != ''">-->
<!--                    {{ clientErrorMessage }}-->
<!--                  </div>-->

<!--                  <div class="form-footer">-->
<!--                    <button v-on:click="addClient" class="btn btn-primary w-100">Add</button>-->
<!--                  </div>-->
<!--                </div>-->
<!--              </div>-->
<!--            </div>-->
<!--            <div class="col-4">-->
<!--              <div class="card">-->
<!--                <div class="card-header">-->
<!--                  <h3 class="card-title">Add new client</h3>-->
<!--                  <div class="card-actions">-->
<!--                    <a href="#" class="btn btn-primary">-->
<!--                      &lt;!&ndash; Download SVG icon from http://tabler-icons.io/i/plus &ndash;&gt;-->
<!--                      <svg xmlns="http://www.w3.org/2000/svg" class="icon" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M12 5l0 14" /><path d="M5 12l14 0" /></svg>-->
<!--                      Add new-->
<!--                    </a>-->
<!--                  </div>-->
<!--                </div>-->
<!--                <div class="card-body p-0">-->
<!--                  <svg xmlns="http://www.w3.org/2000/svg" class="w-100" preserveAspectRatio="none" width="400" height="200" viewBox="0 0 400 200" fill="transparent" stroke="var(&#45;&#45;tblr-border-color, #b8cef1)">-->
<!--                    <line x1="0" y1="0" x2="400" y2="200"></line>-->
<!--                    <line x1="0" y1="200" x2="400" y2="0"></line>-->
<!--                  </svg>-->
<!--                </div>-->
<!--              </div>-->
<!--            </div>-->
<!--            <div class="col-4">-->
<!--              <div class="card">-->
<!--                <div class="card-body" style="height: 10rem"></div>-->
<!--              </div>-->
<!--            </div>-->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>