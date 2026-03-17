<script lang="ts">
  import { onMount } from "svelte";
  import { profilesApi } from "$lib/api";
  import type { Profile } from "$lib/types";
  import { activeProfile } from "$lib/stores/profile";
  import { toast } from "$lib/stores/toast";
  import ProfileForm from "../components/profile/ProfileForm.svelte";
  import ProfileList from "../components/profile/ProfileList.svelte";

  let profileList: Profile[] = $state([]);

  onMount(async () => {
    profileList = await profilesApi.list();
  });

  async function create(name: string) {
    const p = await profilesApi.create(name);
    profileList = [...profileList, p];
    toast.success(`Profile "${p.name}" created`);
  }

  async function del(id: number) {
    await profilesApi.delete(id);
    profileList = profileList.filter((p) => p.id !== id);
    if ($activeProfile?.id === id) activeProfile.set(null);
    toast.info("Profile deleted");
  }

  function select(p: Profile) {
    activeProfile.set(p);
    toast.success(`Switched to ${p.name}`);
  }
</script>

<div class="max-w-xl mx-auto py-8 px-4 space-y-6">
  <h1 class="text-2xl font-bold text-gray-900">Profiles</h1>
  <div class="bg-white border border-gray-200 rounded-xl p-4">
    <p class="text-sm font-medium text-gray-700 mb-3">New profile</p>
    <ProfileForm oncreate={create} />
  </div>
  <ProfileList profiles={profileList} onDelete={del} onSelect={select} />
</div>
