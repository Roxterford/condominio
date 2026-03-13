<script lang="ts">
	import { TipoDeCuota } from '$lib/schemas';
	import {
		Button,
		Checkbox,
		Datepicker,
		Helper,
		Input,
		Label,
		Radio,
		Select,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';
	import { CheckOutline, FileLinesOutline, PlusOutline } from 'flowbite-svelte-icons';
	import { superForm } from 'sveltekit-superforms';
	import type { PageProps } from './$types';
	import AgregarGastoModal from './components/AgregarGastoModal.svelte';


  let { data }: PageProps = $props();

  const { form } = superForm(data.form, {dataType: 'json'});

  let agregar_gasto_modal: ReturnType<typeof AgregarGastoModal>;
</script>

<h1>Crear Nueva Cuota</h1>
<p class="text-gray-500">Complete el formulario para crear una nueva cuota en el sistema</p>
<pre>{JSON.stringify($form, null, 2)}</pre>
<form action="">
 <section>
  <h3>Tipo de Cuota</h3>
  <div class="grid gap-5">
    <div class="radio">a
      <Radio name="tipo" value="REGULAR" bind:group={$form.tipo}>Mensualidad Regular</Radio>
      <Helper class="radio__helper">
        Some helper text here
      </Helper>
    </div>
    <div class="radio">
      <Radio name="tipo" value="ESPECIAL" bind:group={$form.tipo}>Cuota Especial</Radio>
      <Helper class="radio__helper">
        Some helper text here
      </Helper>
    </div>
  </div>
 </section>   
  {#if $form.tipo === TipoDeCuota.Regular}

  <section class="grid gap-5">
  <h3>Periodo</h3>
  <div class="grid grid-cols-2 gap-x-5 gap-y-3">
    <div>
      <Label for="anio">Año</Label>
      <Select id="anio" placeholder="Seleccione el año" disabled={$form.periodo.actual}>
        <option value="2024">2024</option>
        <option value="2025">2025</option>
        <option value="2026">2026</option>
        <option value="2027">2027</option>

      </Select>
    </div>
    <div>
      <Label for="mes">Mes</Label>
      <Select id="mes" placeholder="Seleccione el mes" bind:value={$form.periodo.mes} disabled={$form.periodo.actual}>
        <option value="1">Enero</option>
        <option value="2">Febrero</option>
        <option value="3">Marzo</option>
        <option value="4">Abril</option>
        <option value="5">Mayo</option>
        <option value="6">Junio</option>
        <option value="7">Julio</option>
        <option value="8">Agosto</option>
        <option value="9">Septiembre</option>
        <option value="10">Octubre</option>
        <option value="11">Noviembre</option>
        <option value="12">Diciembre</option>
      </Select>
    </div>

    <div>
      <Checkbox id="useCurrentYear" bind:checked={$form.periodo.actual}>Usar año actual</Checkbox>
    </div>
  </div>
  <div class="grid grid-cols-2 gap-x-5 gap-y-3">
    <div>
      <Label for="fecha_emision">Fecha de Emisión</Label>
      <Datepicker id="fecha_emision" placeholder="dd/mm/aaaa" bind:value={$form.periodo.fecha_emision} />
    </div>
    <div>
      <Label for="fecha_limite">Fecha Límite de Pago</Label>
      <Datepicker id="fecha_limite" placeholder="dd/mm/aaaa" bind:value={$form.periodo.fecha_limite} />
    </div>
  </div>
 </section>  
 <section>
   <div class="flex justify-between items-center">
     <h3>Desglose de Gastos</h3>
     <Button color="alternative" onclick={() => agregar_gasto_modal?.show()}><PlusOutline /> Agregar Gasto</Button>
   </div>
   <Table>
    <TableHead>
      <TableHeadCell>ID</TableHeadCell>
      <TableHeadCell>Concepto</TableHeadCell>
      <TableHeadCell>Monto</TableHeadCell>
      <TableHeadCell>Fecha</TableHeadCell>
      <TableHeadCell>Proveedor</TableHeadCell>
    </TableHead>
    <TableBody>
      <TableBodyRow>
        <TableBodyCell>bogy5a</TableBodyCell>
        <TableBodyCell>Pago de globos</TableBodyCell>
        <TableBodyCell>$ 1,34</TableBodyCell>
        <TableBodyCell>14/12/2025</TableBodyCell>
        <TableBodyCell>Gas Natural</TableBodyCell>
      </TableBodyRow>
   
    </TableBody>
   </Table>
 </section>
 <section>
  <h3>Adjuntar Documentos (opcional)</h3>
  <div class="bg-gray-100 min-h-48 rounded-xl border border-gray-300 border-dashed" style="">
    
  </div>
 </section>
{:else if $form.tipo === TipoDeCuota.Especial}
<section>
  <Label for="descripcion">Descripción</Label>
  <Input type="text" id="descripcion" placeholder="Ej. Reparación de Bomba de Agua" required />
</section>
  {/if}

 <div class="flex justify-end gap-3">
  <Button color="alternative">Cancelar</Button>
  <Button color="alternative"><FileLinesOutline />Guardar Borrador</Button>
  <Button><CheckOutline /> Publicar Cuota</Button>
 </div>
</form>

<AgregarGastoModal bind:this={agregar_gasto_modal}/>
