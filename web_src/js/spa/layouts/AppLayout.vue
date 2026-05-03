<template>
  <div class="page-content">
    <slot/>
  </div>
</template>

<script setup lang="ts">
// AppLayout wraps page content area below the server-rendered navbar.
// The outer navbar, header, and footer are still rendered by the Go template (spa.tmpl),
// so this layout only needs to wrap the inner content.
</script>
