<script>
import CommonMixins from '../mixins/CommonMixins'
import { mailbox } from '../stores/mailbox'

export default {
	mixins: [CommonMixins],

	data() {
		return {
			mailbox,
		}
	},

	methods: {
		inSearch: function (tag) {
			const urlParams = new URLSearchParams(window.location.search)
			const query = urlParams.get('q')
			if (!query) {
				return false
			}

			let re = new RegExp(`(^|\\s)tag:"?${tag}"?($|\\s)`, 'i')
			return query.match(re)
		}
	}
}
</script>

<template>
	<template v-if="mailbox.tags && mailbox.tags.length">
		<div class="mt-4 text-muted">
			<button class="btn btn-sm dropdown-toggle ms-n1" data-bs-toggle="dropdown" aria-expanded="false">
				Tags
			</button>
			<ul class="dropdown-menu dropdown-menu-end">
				<li>
					<button class="dropdown-item" @click="mailbox.showTagColors = !mailbox.showTagColors">
						<template v-if="mailbox.showTagColors">Hide</template>
						<template v-else>Show</template>
						tag colors
					</button>
				</li>
			</ul>
		</div>
		<div class="list-group mt-1 mb-5 pb-3">
			<RouterLink v-for="tag in mailbox.tags" :to="'/search?q=' + tagEncodeURI(tag)" @click="hideNav"
				:style="mailbox.showTagColors ? { borderLeftColor: colorHash(tag), borderLeftWidth: '4px' } : ''"
				class="list-group-item list-group-item-action small px-2" :class="inSearch(tag) ? 'active' : ''">
				<i class="bi bi-tag-fill" v-if="inSearch(tag)"></i>
				<i class="bi bi-tag" v-else></i>
				{{ tag }}
			</RouterLink>
		</div>
	</template>
</template>
