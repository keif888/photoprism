<template>
  <v-snackbar id="p-notify" v-model="visible" :color="message.color" :timeout="-1" variant="elevated" location="bottom">
    {{ message.text }}
    <template #actions>
      <v-btn icon="mdi-close" :color="'on-' + message.color" variant="text" @click="close"></v-btn>
    </template>
  </v-snackbar>
</template>
<script>
import Event from "pubsub-js";

export default {
  name: "PNotify",
  data() {
    return {
      visible: false,
      message: {
        text: "",
        color: "transparent",
      },
      messages: [],
      lastText: "",
      lastId: 1,
      subscriptionId: "",
      defaultColor: "info",
    };
  },
  created() {
    this.subscriptionId = Event.subscribe("notify", this.onNotify);
  },
  unmounted() {
    Event.unsubscribe(this.subscriptionId);
  },
  methods: {
    onNotify: function (ev, data) {
      const type = ev.split(".")[1];

      // Get the message.
      let m = data.message;

      // Skip empty messages.
      if (!m || !m.length) {
        console.warn("notify: empty message");
        return;
      }

      // Log notifications in test mode.
      if (this.$config.test) {
        console.log(type + ": " + m.toLowerCase());
        return;
      }

      // First letter of the message should be uppercase.
      m = m.replace(/^./, m[0].toUpperCase());

      switch (type) {
        case "warning":
          this.addWarningMessage(m);
          break;
        case "error":
          this.addErrorMessage(m);
          break;
        case "success":
          this.addSuccessMessage(m);
          break;
        case "info":
          this.addInfoMessage(m);
          break;
        default:
          alert(m);
      }
    },

    addWarningMessage: function (message) {
      this.addMessage("warning", message, 3000);
    },

    addErrorMessage: function (message) {
      this.addMessage("error", message, 8000);
    },

    addSuccessMessage: function (message) {
      this.addMessage("success", message, 2000);
    },

    addInfoMessage: function (message) {
      this.addMessage("info", message, 2000);
    },

    addMessage: function (color, text, delay) {
      if (text === this.lastText) return;

      this.lastId++;
      this.lastText = text;

      const m = {
        id: this.lastId,
        text,
        color,
        delay,
      };

      this.messages.push(m);

      if (!this.visible) {
        this.show();
      }
    },
    close: function () {
      this.visible = false;
      this.show();
    },
    show: function () {
      const message = this.messages.shift();

      if (message) {
        this.message = message;

        if (!this.message.color) {
          this.message.color = this.defaultColor;
        }

        this.visible = true;

        if (message.delay > 0) {
          setTimeout(() => {
            this.lastText = "";
            this.show();
          }, message.delay);
        }
      } else {
        this.visible = false;
        this.$nextTick(function () {
          this.message.text = "";
          this.message.color = "transparent";
        });
      }
    },
  },
};
</script>
