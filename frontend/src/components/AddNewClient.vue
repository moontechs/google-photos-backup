<template>
  <div class="card">
      <div class="card-header">
        Add new client
      </div>
      <div class="card-body">
          <form method="post">
              <div class="form-group">
                  <p>Client ID</p>
                  <label for="client-id" class="visually-hidden">Client ID</label>
                  <input id="client-id" type="text" class="form-control" required="true">
              </div>

              <div class="form-group">
                  <p>Client Secret</p>
                  <label for="client-secret" class="visually-hidden">Client Secret</label>
                  <input id="client-secret" type="password" class="form-control" required="true">
              </div>
          </form>
      </div>
      <div class="card-footer">
          <div class="row">
              <div class="col-6">
                  <a href="javascript:void(0);" class="btn btn-secondary w-100 _effect--ripple waves-effect waves-light">Add</a>
              </div>
          </div>
      </div>
  </div>
</template>

<script>
export default {
  name: 'AddNewClient',
}
</script>

<style scoped>
</style>
