<script>
import CommonMixins from '../mixins/CommonMixins'
import { pagination } from '../stores/pagination'

export default {
	mixins: [CommonMixins],

	data() {
		return {
			search: '',
			pagination
		}
	},

	mounted() {
		this.searchFromURL()
	},

	watch: {
		$route() {
			this.searchFromURL()
		}
	},

	methods: {
		searchFromURL: function () {
			const urlParams = new URLSearchParams(window.location.search)
			this.search = urlParams.get('q') ? urlParams.get('q') : ''
		},

		doSearch: function (e) {
			pagination.start = 0
			if (this.search == '') {
				this.$router.push('/')
			} else {
				this.$router.push('/search?q=' + encodeURIComponent(this.search))
			}

			e.preventDefault()
		},

		resetSearch: function () {
			this.search = ''
			this.$router.push('/')
		}
	}
}
</script>

<template>
	<form v-on:submit="doSearch">
		<div class="input-group flex-nowrap">
			<div class="ms-md-2 d-flex border bg-body rounded-start flex-fill position-relative">
				<input type="text" class="form-control border-0" aria-label="Search" v-model.trim="search"
					placeholder="Search mailbox">
				<span class="btn btn-link position-absolute end-0 text-muted" v-if="search != ''"
					v-on:click="resetSearch"><i class="bi bi-x-circle"></i></span>
			</div>
			<button class="btn btn-outline-secondary" type="submit">
				<i class="bi bi-search"></i>
			</button>
		</div>
	</form>
</template>
