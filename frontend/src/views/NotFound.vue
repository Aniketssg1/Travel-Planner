<template>
  <div class="notfound-page">
    <div class="notfound-content">
      <div class="notfound-code">404</div>
      <h1 class="notfound-title">Lost in Transit</h1>
      <p class="notfound-sub">
        Looks like this page took a wrong turn somewhere over the Pacific.<br />
        Let's get you back on track.
      </p>
      <RouterLink to="/" class="btn btn-primary btn-arrow btn-lg">
        Back to Explore
        <span class="btn-arrow-icon">↗</span>
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.notfound-page {
  min-height: calc(100vh - var(--nav-h));
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 2rem;
}

.notfound-content {
  max-width: 500px;
}

.notfound-code {
  font-family: var(--font-display);
  font-size: clamp(6rem, 20vw, 12rem);
  font-weight: 900;
  line-height: 1;
  color: var(--bg-warm);
  text-shadow: 0 2px 0 rgba(0, 0, 0, 0.04);
  margin-bottom: 0.5rem;
}

.notfound-title {
  font-size: clamp(1.5rem, 4vw, 2.5rem);
  color: var(--text-dark);
  margin-bottom: 1rem;
}

.notfound-sub {
  color: var(--text-muted);
  font-size: 1rem;
  line-height: 1.7;
  margin-bottom: 2rem;
}
</style>
