import type { FileStorage } from '$lib/domain/ports/FileStorage';
import { SupabaseClient } from '@supabase/supabase-js';

export class SupabaseFileStorage implements FileStorage {
    private supabase: SupabaseClient;
    private bucketName: string;

    constructor(supabaseClient: SupabaseClient, bucketName: string = 'uploads') {
        this.supabase = supabaseClient;
        this.bucketName = bucketName;
    }

    /**
     * Sube un archivo a un bucket y carpeta específicos.
     * @returns La URL pública del archivo subido.
     */
    async upload(file: File, store: string): Promise<string> {
        const filePath = `${store}/${Date.now()}_${file.name}`;

        const { data, error } = await this.supabase.storage
            .from(this.bucketName)
            .upload(filePath, file);

        if (error) {
            throw new Error(`Error subiendo archivo: ${error.message}`);
        }

        // Obtenemos la URL pública para cumplir con el contrato de devolver un string
        const { data: { publicUrl } } = this.supabase.storage
            .from(this.bucketName)
            .getPublicUrl(data.path);

        return publicUrl;
    }

    /**
     * Elimina un archivo dado su path relativo dentro del bucket.
     */
    async delete(path: string): Promise<void> {
        const { error } = await this.supabase.storage
            .from(this.bucketName)
            .remove([path]);

        if (error) {
            throw new Error(`Error eliminando archivo: ${error.message}`);
        }
    }
}