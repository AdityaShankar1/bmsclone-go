<script>
  import { onMount } from 'svelte';
  import MovieCard from './MovieCard.svelte';
  import Bookings from './Bookings.svelte';
  import TheatreSelector from './TheatreSelector.svelte';
  import Wallet from './Wallet.svelte';

  const API_BASE = "http://localhost:8080";
  
  let view = 'home';
  let bookingFlow = 'movie';
  let checkoutStep = 'select'; // select, pay, confirm

  let movies = [];
  let filteredMovies = [];
  let genreFilter = '';
  let formatFilter = '';
  
  let selectedMovie = null;
  let selectedTheatre = null;
  let selectedFormat = null;
  let selectedShow = null;
  let selectedSeat = null;
  
  let shows = [];
  let seats = [];
  let ticket = null;
  let paymentMethod = 'wallet'; // wallet, upi, card
  let loading = false;

  async function fetchMovies() {
    loading = true;
    try {
      const res = await fetch(`${API_BASE}/movies`);
      const data = await res.json();
      movies = data.movies || [];
      applyFilters();
    } catch (err) { console.error(err); }
    finally { loading = false; }
  }

  function applyFilters() {
    filteredMovies = movies.filter(m => {
      const matchesGenre = genreFilter === '' || m.genre === genreFilter;
      const matchesFormat = formatFilter === '' || (m.formats && m.formats.includes(formatFilter));
      return matchesGenre && matchesFormat;
    });
  }

  async function fetchShows() {
    loading = true;
    let url = `${API_BASE}/shows?`;
    if (selectedMovie) url += `movie_id=${selectedMovie.id}&`;
    if (selectedTheatre) url += `theatre_id=${selectedTheatre.id}&`;
    const res = await fetch(url);
    const data = await res.json();
    shows = data.shows || [];
    loading = false;
  }

  async function fetchSeats() {
    loading = true;
    const res = await fetch(`${API_BASE}/seats?show_id=${selectedShow.id}`);
    const data = await res.json();
    seats = data.seats || [];
    loading = false;
  }

  const handleMovieSelect = (m) => { selectedMovie = m; bookingFlow = 'movie'; view = 'formats'; };
  const handleFormatSelect = (f) => { selectedFormat = f; view = 'theatres'; };
  const handleTheatreSelect = (t) => { selectedTheatre = t; view = 'shows'; fetchShows(); };
  const handleShowSelect = (s) => { selectedShow = s; view = 'seats'; fetchSeats(); };

  async function lockSeat(seat) {
    if (seat.status !== 'available') return;
    try {
      const res = await fetch(`${API_BASE}/book/${seat.id}?show_id=${selectedShow.id}`, { method: 'POST' });
      if (res.status === 200) { selectedSeat = seat.id; }
      else { const d = await res.json(); alert(d.error); }
    } catch (err) { alert("Lock failed"); }
  }

  async function confirmBooking() {
    loading = true;
    try {
      const res = await fetch(`${API_BASE}/confirm-booking`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ seat_id: selectedSeat, show_id: selectedShow.id })
      });
      const data = await res.json();
      if (res.status === 200) { ticket = data; view = 'ticket'; }
      else { alert(data.error); }
    } catch (err) { alert("Payment error"); }
    finally { loading = false; }
  }

  onMount(fetchMovies);
  $: { genreFilter; formatFilter; applyFilters(); }
</script>

<nav class="navbar navbar-expand navbar-dark bg-danger mb-4 shadow sticky-top">
  <div class="container">
    <span class="navbar-brand fw-bold cursor-pointer" on:click={() => { view = 'home'; ticket = null; selectedMovie = null; selectedTheatre = null; }}>BookMyShow</span>
    <div class="navbar-nav ms-auto gap-2">
      <button class="btn btn-outline-light btn-sm" on:click={() => view = 'wallet'}>Wallet</button>
      <button class="btn btn-outline-light btn-sm" on:click={() => view = 'history'}>My Bookings</button>
    </div>
  </div>
</nav>

<main class="container mb-5 text-dark">
  {#if view === 'home'}
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2 class="fw-bold m-0">Recommended Movies</h2>
      <button class="btn btn-dark" on:click={() => { bookingFlow = 'cinema'; view = 'theatres'; selectedMovie = null; }}>Browse by Cinema</button>
    </div>

    <div class="row mb-4 align-items-end g-3">
      <div class="col-md-4">
        <label class="form-label small text-muted fw-bold">GENRE</label>
        <select class="form-select border-0 shadow-sm" bind:value={genreFilter}>
          <option value="">All Genres</option>
          <option value="Action">Action</option>
          <option value="Adventure">Adventure</option>
          <option value="Sci-Fi">Sci-Fi</option>
        </select>
      </div>
      <div class="col-md-4">
        <label class="form-label small text-muted fw-bold">FORMAT</label>
        <select class="form-select border-0 shadow-sm" bind:value={formatFilter}>
          <option value="">All Formats</option>
          <option value="IMAX">IMAX</option>
          <option value="2D">2D</option>
          <option value="3D">3D</option>
          <option value="4DX">4DX</option>
        </select>
      </div>
    </div>

    {#if loading && movies.length === 0}
      <div class="text-center py-5"><div class="spinner-border text-danger"></div></div>
    {:else}
      <div class="row g-4">
        {#each filteredMovies as movie}
          <div class="col-6 col-md-4 col-lg-3">
            <MovieCard {movie} onSelect={handleMovieSelect} />
          </div>
        {/each}
      </div>
    {/if}

  {:else if view === 'formats'}
    <div class="text-center py-4">
      <button class="btn btn-link text-decoration-none text-muted mb-3" on:click={() => view = 'home'}>← Back</button>
      <h2 class="fw-bold mb-4">Select Format for <span class="text-danger">{selectedMovie.title}</span></h2>
      <div class="d-flex flex-wrap justify-content-center gap-3">
        {#each selectedMovie.formats.split(',') as format}
          <button class="btn btn-outline-danger btn-lg px-5 py-3 shadow-sm fw-bold" on:click={() => handleFormatSelect(format)}>{format}</button>
        {/each}
      </div>
    </div>

  {:else if view === 'theatres'}
    <div class="py-2">
      <button class="btn btn-link text-decoration-none text-muted mb-3" on:click={() => view = 'home'}>← Back</button>
      <TheatreSelector 
        movieFilter={bookingFlow === 'movie' ? selectedMovie : null} 
        formatFilter={bookingFlow === 'movie' ? selectedFormat : null}
        on:select={(e) => handleTheatreSelect(e.detail)} />
    </div>

  {:else if view === 'shows'}
    <div class="py-4">
      <button class="btn btn-link text-decoration-none text-muted mb-3" on:click={() => view = 'theatres'}>← Back</button>
      <h3 class="fw-bold mb-1">{selectedTheatre.name}</h3>
      <div class="row g-3 mt-3">
        {#each shows as show}
          <div class="col-md-6 col-lg-4">
            <button class="card w-100 text-start border-0 shadow-sm p-3 show-card" on:click={() => handleShowSelect(show)}>
              <span class="badge bg-danger mb-2 align-self-start">{show.format}</span>
              <h5 class="fw-bold mb-1">{movies.find(m => m.id === show.movie_id)?.title}</h5>
              <h6 class="text-success fw-bold">₹{show.price}</h6>
              <p class="mb-0 text-muted">{new Date(show.showtime).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}</p>
            </button>
          </div>
        {/each}
      </div>
    </div>

  {:else if view === 'seats'}
    <div class="container py-2 text-center">
      <div class="d-flex justify-content-between align-items-center mb-5">
        <button class="btn btn-outline-secondary" on:click={() => view = 'shows'}>← Back</button>
        <h4 class="fw-bold mb-0">{selectedTheatre.name}</h4>
      </div>

      <div class="row g-2 justify-content-center mx-auto mb-5" style="max-width: 500px;">
        {#each seats as seat}
          <div class="col-3">
            <button 
              class="btn btn-sm w-100 py-3 {selectedSeat === seat.id ? 'btn-primary active' : (seat.status === 'booked' ? 'btn-danger' : 'btn-outline-secondary')}"
              disabled={seat.status !== 'available' && selectedSeat !== seat.id}
              on:click={() => lockSeat(seat)}>{seat.id}</button>
          </div>
        {/each}
      </div>

      {#if selectedSeat}
        <div class="shadow-lg p-4 rounded bg-white border mx-auto mb-5" style="max-width: 400px;">
          <h5 class="fw-bold mb-4">Payment Interface</h5>
          <div class="list-group mb-4">
            <button class="list-group-item list-group-item-action d-flex align-items-center gap-3 py-3 {paymentMethod === 'wallet' ? 'active' : ''}" on:click={() => paymentMethod = 'wallet'}>
              <i class="bi bi-wallet2 h4 mb-0"></i>
              <div class="text-start">
                 <span class="d-block fw-bold">Wallet / BMSCash</span>
                 <span class="small opacity-75">Pay using credits & investments</span>
              </div>
            </button>
            <button class="list-group-item list-group-item-action d-flex align-items-center gap-3 py-3 {paymentMethod === 'upi' ? 'active' : ''}" on:click={() => paymentMethod = 'upi'}>
              <i class="bi bi-qr-code h4 mb-0"></i>
              <div class="text-start">
                 <span class="d-block fw-bold">UPI</span>
                 <span class="small opacity-75">GPay, PhonePe, Paytm</span>
              </div>
            </button>
            <button class="list-group-item list-group-item-action d-flex align-items-center gap-3 py-3 {paymentMethod === 'card' ? 'active' : ''}" on:click={() => paymentMethod = 'card'}>
              <i class="bi bi-credit-card h4 mb-0"></i>
              <div class="text-start">
                 <span class="d-block fw-bold">Credit/Debit Card</span>
                 <span class="small opacity-75">All major cards supported</span>
              </div>
            </button>
          </div>
          <button class="btn btn-danger btn-lg w-100" on:click={confirmBooking}>Pay ₹{selectedShow.price}</button>
        </div>
      {/if}
    </div>

  {:else if view === 'ticket'}
    <div class="card mx-auto mt-5 shadow-lg border-danger" style="max-width: 400px;">
      <div class="card-header bg-danger text-white text-center py-3"><h3>Confirmed! 🎟️</h3></div>
      <div class="card-body text-center p-4">
        <h2 class="text-danger fw-bold mb-4">{selectedSeat}</h2>
        <p class="mb-4">Ticket details sent to your mobile.</p>
        <button class="btn btn-outline-danger w-100" on:click={() => { view = 'home'; ticket = null; selectedSeat = null; }}>Back to Home</button>
      </div>
    </div>

  {:else if view === 'history'}
    <Bookings onBack={() => view = 'home'} />
  {:else if view === 'wallet'}
    <div class="container py-2">
      <button class="btn btn-link text-decoration-none text-muted mb-3" on:click={() => view = 'home'}>← Back</button>
      <Wallet />
    </div>
  {/if}
</main>

<style>
  .cursor-pointer { cursor: pointer; }
  .show-card { transition: all 0.2s; cursor: pointer; }
  .show-card:hover { transform: translateY(-3px); box-shadow: 0 5px 15px rgba(0,0,0,0.1) !important; background: #fff5f5; }
</style>
