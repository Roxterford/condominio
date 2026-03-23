export interface FileStorage {
    upload(file: File, store: string): Promise<string>;
    delete(path: string): Promise<void>;
}