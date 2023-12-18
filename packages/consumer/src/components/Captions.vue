<template>
  <div>
    <div id="first" class="caption">{{ first }}</div>
    <div id="second" class="caption">{{ second }}</div>
  </div>
</template>
<script>
export default {
  props: {
    id: String,
  },
  data() {
    return {
      buffer: [],
      delay: 25,
      first: "Laddar undertexter...",
      second: "...",
    }
  },
  computed: {
  },
  mounted() {
    this.flexFont();
    window.addEventListener("resize", this.flexFont);
    setInterval(() => { this.GETlivecap() }, 150)
    setTimeout(() => { this.renderCaps() }, this.delay)
  },
  beforeDestroy() {
  },
  methods: {
    GETlivecap() {
			fetch("https://sttylus.se/ws/caption/" + this.id)
			.then(resp => resp.json())
			.then(data =>
				{
          this.buffer.push(data)
        }
			)
			.catch(err => {console.log(err)})
    },
    renderCaps() {
      this.first = ""
      this.second = ""
      if (this.buffer.length > 0) {
        let block = this.buffer.shift()
        if (block.Line1) {
          this.first = block.Line1
        }
        if (block.Line2) {
          this.second = block.Line2
        }
      }
      setTimeout(() =>  { this.renderCaps() }, 150)
    },
    flexFont() {
      var lines = document.getElementsByClassName("caption");
      for(var i = 0; i < lines.length; i++) {
        var relFontsize = lines[i].offsetWidth*0.05*0.92;
        lines[i].style.fontSize = relFontsize+'px';
      }
    }
  }
}
</script>
<style>
.caption {
  font-family: 'Roboto Mono', monospace;
  background-color: #000;

  width: 80%;
  color: #fff;
  margin: auto;
  margin-top: 0.5rem;
}
</style>
