<script>
  import { onMount, createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
  const API_BASE = "http://localhost:8080";

  export let movieFilter = null;
  export let formatFilter = null;

  let chains = [];
  let selectedChain = null;
  let theatres = [];
  let loading = false;

  async function fetchChains() {
    const res = await fetch(`${API_BASE}/chains`);
    const data = await res.json();
    chains = data.chains || [];
  }

  async function fetchTheatres(chainID = null) {
    loading = true;
    let url = `${API_BASE}/theatres?`;
    if (chainID) url += `chain_id=${chainID}&`;
    if (movieFilter) url += `movie_id=${movieFilter.id}&`;
    if (formatFilter) url += `format=${formatFilter}&`;

    const res = await fetch(url);
    const data = await res.json();
    theatres = data.theatres || [];
    loading = false;
  }

  function handleChainSelect(chain) {
    selectedChain = chain;
    fetchTheatres(chain.id);
  }

  onMount(() => {
    fetchChains();
    if (movieFilter) fetchTheatres();
  });
</script>

<div class="row g-4">
  {#if !movieFilter}
    <div class="col-12">
      <h5 class="fw-bold mb-3">Filter by Chain</h5>
      <div class="d-flex flex-wrap gap-2">
        <button class="btn {!selectedChain ? 'btn-danger' : 'btn-outline-danger'}" on:click={() => { selectedChain = null; fetchTheatres(); }}>All Chains</button>
        {#each chains as chain}
          <button class="btn {selectedChain?.id === chain.id ? 'btn-danger' : 'btn-outline-danger'}" on:click={() => handleChainSelect(chain)}>
            {chain.name}
          </button>
        {/each}
      </div>
    </div>
  {/if}

  <div class="col-12">
    <h5 class="fw-bold mb-3">{movieFilter ? `Theatres showing ${movieFilter.title}` : 'Select a Theatre'}</h5>
    {#if loading}
      <div class="spinner-border text-danger"></div>
    {:else}
      <div class="list-group shadow-sm border-0">
        {#each theatres as theatre}
          {@const supportsFormat = !formatFilter || theatre.supported_formats.includes(formatFilter)}
          <button 
            class="list-group-item list-group-item-action d-flex justify-content-between align-items-center py-3 {!supportsFormat ? 'opacity-50 grayscale' : ''}"
            disabled={!supportsFormat}
            on:click={() => dispatch('select', theatre)}>
            <div>
              <h6 class="mb-1 fw-bold">{theatre.name}</h6>
              <p class="mb-0 text-muted small">{theatre.address || 'Bengaluru'}</p>
            </div>
            <div class="text-end">
              <div class="d-flex gap-1 mb-1 justify-content-end">
                {#each theatre.supported_formats.split(',') as f}
                  <span class="badge {f === formatFilter ? 'bg-danger' : 'bg-light text-dark border'}">{f}</span>
                {/each}
              </div>
              {#if !supportsFormat}
                <span class="text-danger small fw-bold">Format unavailable</span>
              {/if}
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .grayscale { filter: grayscale(1); }
</style>
