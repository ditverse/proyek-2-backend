-- Migration: Add ON DELETE behavior to organisasi foreign keys
-- Purpose: Allow deletion of organisasi records by setting proper ON DELETE behavior

-- First, drop existing foreign key constraint on users table
ALTER TABLE public.users
DROP CONSTRAINT IF EXISTS users_organisasi_kode_fkey;

-- Recreate the foreign key with ON DELETE SET NULL
-- This will set organisasi_kode to NULL when the referenced organisasi is deleted
ALTER TABLE public.users
ADD CONSTRAINT users_organisasi_kode_fkey 
FOREIGN KEY (organisasi_kode) 
REFERENCES public.organisasi(kode_organisasi) 
ON DELETE SET NULL;

-- Also update kegiatan table's foreign key to organisasi
ALTER TABLE public.kegiatan
DROP CONSTRAINT IF EXISTS kegiatan_organisasi_kode_fkey;

ALTER TABLE public.kegiatan
ADD CONSTRAINT kegiatan_organisasi_kode_fkey 
FOREIGN KEY (organisasi_kode) 
REFERENCES public.organisasi(kode_organisasi) 
ON DELETE SET NULL;
