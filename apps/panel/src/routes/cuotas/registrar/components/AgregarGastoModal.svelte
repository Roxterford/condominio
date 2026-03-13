<script lang="ts">
	import { Modal, Radio } from "flowbite-svelte";

    import { Badge, Button } from "flowbite-svelte";

  let open = $state(true);

  export const show = () => (open = true);
  export const hide = () => (open = false);

</script>

<style>
  @reference "tailwindcss";

  .radiolist {
    @apply grid gap-5;
  }
  
  .radiolist__item {
    @apply flex gap-2 justify-between border rounded-xl p-5 items-start;
  }
  
  .radiolist__item__label {
    @apply font-semibold;
  }
  
  .radiolist__item__description {
    @apply text-sm;
    color: var(--color-muted);

  }
</style>

<Modal bind:open title="Agregar Gasto" class="text-font">
<form class="grid gap-5" method="dialog">
  <ul class="radiolist">
    <li class="radiolist__item-content">
      <label class="radiolist__item" for="no_asociados">
        <div>
          <span class="radiolist__item__label">Seleccionar gastos no asociados</span>
          <p class="radiolist__item__description">Vincule gastos huerfanos registrados previamente</p>
        </div>
        <Badge color="green">Recomendado</Badge>
        <Radio name="agregar_gasto" id="no_asociados" defaultChecked/>
      </label>
    </li>
    <li class="radiolist__item-content">
      <label class="radiolist__item" for="todos_asociados">
        <div>
          <span class="radiolist__item__label">Añadir todos los gastos</span>
          <p class="radiolist__item__description">Añade todos los gastos huerfanos automaticamente</p>
        </div>
        <Radio type="radio" name="agregar_gasto" id="todos_asociados" />
      </label>
    </li>
    <li class="radiolist__item-content">
      <label class="radiolist__item" for="nuevo_gasto">
        <div>
          <span class="radiolist__item__label">Registrar nuevo gasto</span>
          <p class="radiolist__item__description">Registre un nuevo gasto en el sistema</p>
        </div>
        <Radio type="radio" name="agregar_gasto" id="nuevo_gasto" />
      </label>
    </li>
  </ul>
  <div class="flex justify-end gap-2">
    <Button type="button" color="alternative" onclick={hide}>Cancelar</Button>
    <Button type="submit" color="primary">Siguiente</Button>
  </div>
</form>
</Modal>