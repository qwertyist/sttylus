<template>
  <div>
    <b-form
      v-if="show"
      @submit="onSubmit"
      @reset="onReset"
    >
      <b-form-group
        id="user-settings-fullname"
        label="För- och efternamn:"
        label-for="user-fullname-input"
        description="Vi delar inte din information med någon utomstående"
      >
        <b-form-input
          v-model="form.name"
          id="user-fullname-input"
          disabled
          placeholder="Skriv ditt namn..."
        />
      </b-form-group>
      <b-form-group
        id="user-settings-email"
        label="E-postadress:"
        label-for="user-fullname-input"
      >
        <b-form-input
          v-model="form.email"
          id="user-email-input"
          type="email"
          disabled
          placeholder="Skriv din e-postadress"
        />
      </b-form-group>
      <b-form-group
        id="user-settings-phone"
        label="Telefonnummer:"
        label-for="user-phone-input"
      >
        <b-form-input
          v-model="form.phone"
          id="user-phone-input"
          type="tel"
          disabled
          placeholder="Skriv ditt telefonnummer"
        />
      </b-form-group>
      <!--<b-button type="submit" variant="primary">Skicka in</b-button>
      <b-button type="reset" variant="danger">Rensa formulär</b-button>-->
    </b-form>
  </div>
</template>

<script>
export default {
  name: 'User',
  data () {
    return {
      form: {
        name: '',
        company: '',
        description: '',
        email: '',
        phone: '',
        role: ''
      },
      show: true
    };
  },
  methods: {
    onSubmit (evt) {
      evt.preventDefault();
      const data = JSON.stringify(this.form);
      localStorage.setItem('user', data);
    },
    onReset (evt) {
      evt.preventDefault();

      this.form.name = '';
      this.form.email = '';
      this.form.phone = '';
      this.show = false;
      this.$nextTick(() => {
        this.show = true;
      });
    }
  },
  mounted () {
    this.form = this.$store.state.userData;
   
  }
};
</script>