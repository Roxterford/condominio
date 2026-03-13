import { NuevaCuetoSchema, TipoDeCuota } from '$lib/schemas';
import { fail } from '@sveltejs/kit';
import { superValidate } from 'sveltekit-superforms';
import { valibot } from 'sveltekit-superforms/adapters';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
  const form = await superValidate(valibot(NuevaCuetoSchema), {
    defaults: {
      tipo: TipoDeCuota.Regular,
      periodo: {
        actual: true,
        anio: new Date().getFullYear(),
        mes: new Date().getMonth() + 1,
        fecha_emision: new Date(),
        fecha_limite: new Date()
      }
    }
  });
  return { form };
};

export const actions = {
  default: async ({ request }) => {
    // Superforms se encarga de parsear el JSON y los archivos automáticamente
    const form = await superValidate(request, valibot(NuevaCuetoSchema) );

    if (!form.valid) return fail(400, { form });


    return { form };
  }
};