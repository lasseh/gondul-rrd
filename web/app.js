const distroTree = "https://nms.tg18.gathering.org/api/public/distro-tree";
const switchState = "https://nms.tg18.gathering.org/api/public/switch-state";

const vm = new Vue({
	el: '#app',
	data: {
		results: [],
		selected: undefined,
		totals: undefined,
		interfaces: [],
		totalsImg: undefined,
		defaultStart: "8h",
		defaultEnd: "-1m",
		rrdImages: [],
		device: undefined,
	},
	mounted() {
		axios.get(distroTree).then(response => {
			this.results = response.data["distro-tree"]
		})
	},
	methods: {
		getSwitch: function() {
			console.log(this.selected)
			this.getInterfaces()
		},
		getInterfaces: function(selected) {
			//console.log("Getting interfaces for: " + this.selected)
			axios.get(switchState + "/" + this.selected).then(response => {
				this.interfaces = response.data["switches"][this.selected]
			})
			this.totalsImg = 'https://rrd.lasse.cloud/graph?width=700&height=200&legend=1&start=-80h&end=-1h&device=' + this.selected + '&interface=totals&title=Totals'
			//console.log(this.interfaces)
			console.log(this.getImage(this.selected))

		},
		getImage: function(device, interface, alias) {

			rrdUrl = 'https://rrd.lasse.cloud/graph?width=400&height=100&legend=1&start=-80h&end=-1h&device=' + 
				device + 
				'&interface='+
				interface+ 
				'&title='+
				alias
			return rrdUrl

		},
	}
});
