<script>
  export let movie;
  export let onSelect;
</script>

<div class="card h-100 shadow-sm border-0 movie-card" on:click={() => onSelect(movie)}>
  <div class="position-relative">
    <img src={movie.poster_url} class="card-img-top" alt={movie.title} style="height: 350px; object-fit: cover;">
    <span class="badge bg-dark position-absolute top-0 end-0 m-2 opacity-75">{movie.rating}</span>
  </div>
  <div class="card-body">
    <h5 class="card-title fw-bold text-dark">{movie.title}</h5>
    <p class="text-muted small mb-2">{movie.genre}</p>
    <div class="d-flex flex-wrap gap-1">
      {#each movie.formats.split(',') as format}
        <span class="badge bg-light text-dark border">{format}</span>
      {/each}
    </div>
  </div>
</div>

<style>
  .movie-card {
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s;
  }
  .movie-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0,0,0,0.1) !important;
  }
</style>
