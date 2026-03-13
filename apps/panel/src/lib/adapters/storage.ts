import { supabase } from '$lib/providers/supabase';
import { SupabaseFileStorage } from './SupabaseFileStorage';

// Instancia única (Singleton) lista para usar
export const storage = new SupabaseFileStorage(supabase);