<script>
  import { onMount } from 'svelte';
  const API_BASE = "http://localhost:8080";
  let bookings = [];
  let loading = true;

  export let onBack;

  async function fetchBookings() {
    try {
      const res = await fetch(`${API_BASE}/bookings`);
      const data = await res.json();
      bookings = data.bookings || [];
      loading = false;
    } catch (err) {
      console.error(err);
      loading = false;
    }
  }

  async function handleCancel(id) {
    if (!confirm("Are you sure? This ticket will be cancelled without refund.")) return;
    const res = await fetch(`${API_BASE}/cancel-booking/${id}`, { method: 'POST' });
    if (res.status === 200) { alert("Ticket Cancelled"); fetchBookings(); }
  }

  async function handleSeekRefund(id) {
    if (!confirm("Seek 70% refund to BMSCash?")) return;
    const res = await fetch(`${API_BASE}/cancel-booking/${id}`, { method: 'POST' });
    if (res.status === 200) { 
        alert("Refunded successfully!"); 
        fetchBookings(); 
    }
  }

  async function rateTicket(id, rating) {
    const res = await fetch(`${API_BASE}/rate-ticket/${id}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ rating })
    });
    if (res.status === 200) { fetchBookings(); }
  }

  function canRefund(showtime) {
    const now = new Date();
    const show = new Date(showtime);
    const diffHours = (show - now) / (1000 * 60 * 60);
    return diffHours > 2; // Only refund if > 2 hours away
  }

  onMount(fetchBookings);
</script>

<div class="container py-4 px-3">
  <div class="d-flex justify-content-between align-items-center mb-4">
    <button class="btn btn-outline-secondary" on:click={onBack}>← Back to Movies</button>
    <h2 class="fw-bold m-0 text-dark">My Bookings</h2>
  </div>

  {#if loading}
    <div class="text-center py-5"><div class="spinner-border text-danger"></div></div>
  {:else if bookings.length === 0}
    <div class="text-center py-5 text-dark">
      <h4 class="text-muted">No bookings found yet.</h4>
    </div>
  {:else}
    <div class="row g-4">
      {#each bookings as booking}
        <div class="col-md-6 col-lg-4">
          <div class="card border-0 shadow-sm text-dark h-100 {booking.status === 'Cancelled' ? 'opacity-75' : ''}">
            <div class="card-body">
              <div class="d-flex justify-content-between mb-3">
                <span class="badge {booking.status === 'Confirmed' ? 'bg-success' : 'bg-danger'}">{booking.status}</span>
                <span class="text-muted small">#BMS-{booking.id}</span>
              </div>
              
              <h5 class="fw-bold mb-1">{booking.movie_title}</h5>
              <p class="text-muted small mb-3">{booking.format} • {new Date(booking.showtime).toLocaleString()}</p>
              
              <div class="bg-light p-3 rounded mb-3 d-flex justify-content-between align-items-center">
                <div>
                  <span class="d-block small text-muted">SEAT</span>
                  <span class="fw-bold h4 m-0">{booking.seat_id}</span>
                </div>
                <div class="text-end">
                  <span class="d-block small text-muted">PAID</span>
                  <span class="fw-bold h5 m-0">₹{booking.price_paid.toFixed(2)}</span>
                </div>
              </div>

              {#if booking.status === 'Confirmed'}
                <div class="d-flex gap-2 mb-3">
                   <button class="btn btn-outline-danger btn-sm flex-grow-1" on:click={() => handleCancel(booking.id)}>Cancel</button>
                   <button 
                     class="btn btn-outline-primary btn-sm flex-grow-1" 
                     disabled={!canRefund(booking.showtime)}
                     on:click={() => handleSeekRefund(booking.id)}>
                     Seek Refund
                   </button>
                </div>
                
                <div class="text-center border-top pt-3">
                  <p class="small text-muted mb-2">Rate your movie experience:</p>
                  <div class="d-flex justify-content-center gap-1">
                    {#each [1,2,3,4,5] as star}
                      <button 
                        class="btn btn-sm btn-link p-0 text-decoration-none" 
                        on:click={() => rateTicket(booking.id, star)}>
                        <i class="bi {booking.user_rating >= star ? 'bi-star-fill text-warning' : 'bi-star'} h5"></i>
                      </button>
                    {/each}
                  </div>
                </div>
              {:else}
                <div class="alert alert-secondary py-2 small mb-0 text-center">Ticket Refunded to Wallet</div>
              {/if}
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
