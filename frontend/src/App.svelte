<script lang="ts">
  import Router from "svelte-spa-router";
  import Navbar from "./components/layout/Navbar.svelte";
  import Toast from "./components/layout/Toast.svelte";
  import LibraryPage from "./pages/LibraryPage.svelte";
  import CatalogPage from "./pages/CatalogPage.svelte";
  import BookPage from "./pages/BookPage.svelte";
  import AddBookPage from "./pages/AddBookPage.svelte";
  import ProfilesPage from "./pages/ProfilesPage.svelte";

  const routes = {
    "/": LibraryPage,
    "/catalog": CatalogPage,
    "/books/add": AddBookPage,
    "/books/:id": BookPage,
    "/profiles": ProfilesPage,
  };
</script>

<div class="min-h-screen bg-gray-50">
  <Navbar />
  <main>
    <Router {routes} />
  </main>
  <Toast />
</div>
