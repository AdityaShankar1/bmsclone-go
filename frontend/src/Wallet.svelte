<script>
  import { onMount } from 'svelte';
  const API_BASE = "http://localhost:8080";

  let account = { bms_cash: 0, investment_balance: 0 };
  let transactions = [];
  let investAmount = 500;
  let loading = true;

  async function fetchData() {
    loading = true;
    const [accRes, txRes] = await Promise.all([
      fetch(`${API_BASE}/account`),
      fetch(`${API_BASE}/transactions`)
    ]);
    account = await accRes.json();
    const txData = await txRes.json();
    transactions = txData.transactions || [];
    loading = false;
  }

  async function handleInvest() {
    const res = await fetch(`${API_BASE}/invest`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ amount: Number(investAmount) })
    });
    if (res.status === 200) {
      alert("Amount added to wallet!");
      fetchData();
    }
  }

  onMount(fetchData);
</script>

<div class="container py-4">
  <!-- Large Balance Header -->
  <div class="text-center mb-5">
    <h6 class="text-muted text-uppercase letter-spacing-2 mb-2">Total Wallet Balance</h6>
    <h1 class="display-3 fw-bold text-dark">₹{(account.bms_cash + account.investment_balance).toFixed(2)}</h1>
    <div class="d-flex justify-content-center gap-4 mt-3">
      <div class="text-center">
        <span class="d-block small text-muted">BMSCash</span>
        <span class="fw-bold h5 text-danger">₹{account.bms_cash.toFixed(2)}</span>
      </div>
      <div class="vr"></div>
      <div class="text-center">
        <span class="d-block small text-muted">Invested</span>
        <span class="fw-bold h5 text-primary">₹{account.investment_balance.toFixed(2)}</span>
      </div>
    </div>
  </div>

  <div class="row g-4">
    <!-- Top-up Section -->
    <div class="col-lg-4">
      <div class="card border-0 shadow-lg bg-dark text-white p-4 h-100">
        <h5 class="fw-bold mb-4">Add Money</h5>
        <div class="mb-4">
          <label class="small opacity-75 mb-2">Amount (₹)</label>
          <input type="number" class="form-control form-control-lg bg-transparent text-white border-secondary" bind:value={investAmount}>
        </div>
        <button class="btn btn-danger btn-lg w-100 shadow-sm" on:click={handleInvest}>Top Up Wallet</button>
        
        <div class="mt-4 p-3 bg-secondary bg-opacity-25 rounded">
          {#if account.investment_balance >= 500}
            <div class="text-success small">
              <i class="bi bi-shield-fill-check"></i> Inflation Protection Active
              <p class="mt-1 mb-0 opacity-75">You get 20% discount on all bookings!</p>
            </div>
          {:else}
            <div class="text-warning small">
              <i class="bi bi-shield-fill-exclamation"></i> Invest ₹{500 - account.investment_balance} more to freeze rates.
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Transaction History -->
    <div class="col-lg-8">
      <div class="card border-0 shadow-sm p-4 h-100">
        <h5 class="fw-bold mb-4">Transaction History</h5>
        <div class="table-responsive" style="max-height: 400px;">
          <table class="table table-hover align-middle">
            <thead class="text-muted small text-uppercase">
              <tr>
                <th>Date</th>
                <th>Description</th>
                <th>Type</th>
                <th class="text-end">Amount</th>
              </tr>
            </thead>
            <tbody>
              {#each transactions as tx}
                <tr>
                  <td class="small text-muted">{new Date(tx.created_at).toLocaleDateString()}</td>
                  <td>
                    <span class="d-block fw-bold">{tx.description}</span>
                  </td>
                  <td>
                    <span class="badge {tx.type === 'Payment' ? 'bg-danger-subtle text-danger' : (tx.type === 'Investment' ? 'bg-success-subtle text-success' : 'bg-primary-subtle text-primary')}">
                      {tx.type}
                    </span>
                  </td>
                  <td class="text-end fw-bold">
                    {tx.type === 'Payment' ? '-' : '+'}₹{tx.amount.toFixed(2)}
                  </td>
                </tr>
              {/each}
              {#if transactions.length === 0 && !loading}
                <tr><td colspan="4" class="text-center py-4 text-muted">No transactions yet.</td></tr>
              {/if}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .letter-spacing-2 { letter-spacing: 2px; }
  .vr { width: 1px; background-color: #dee2e6; opacity: 0.3; }
</style>
