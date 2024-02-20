-- Just one step where we clean all data from file_data table

-- Add real null values to string fields
update public.file_data fd set cpf = NULL where fd.cpf = 'NULL';
update public.file_data fd set loja_mais_frequente = NULL where fd.loja_mais_frequente = 'NULL';
update public.file_data fd set loja_ultima_compra = NULL where fd.loja_ultima_compra = 'NULL';

-- Clean CPF/CNPJ with invalid data
delete from public.file_data fd where fd.cpf ~ '[a-zA-Z]+' or length(fd.cpf) < 11;
delete from public.file_data fd where fd.loja_mais_frequente ~ '[a-zA-Z]+' or length(fd.loja_mais_frequente) < 14;
delete from public.file_data fd where fd.loja_ultima_compra ~ '[a-zA-Z]+' or length(fd.loja_ultima_compra) < 14;

-- Corretly format data for CPF/CNPJ columns
update public.file_data fd set cpf = public.format_cpf_cnpj(cpf)
where  not fd.cpf ~ '\y[0-9]{3}.[0-9]{3}.[0-9]{3}-[0-9]{2}\y';

update public.file_data fd set loja_mais_frequente = public.format_cpf_cnpj(loja_mais_frequente)
where not fd.loja_mais_frequente ~ '\y[0-9]{2}.[0-9]{3}.[0-9]{3}/[0-9]{4}-[0-9]{2}\y';

update public.file_data fd set loja_ultima_compra = public.format_cpf_cnpj(loja_ultima_compra)
where not fd.loja_ultima_compra ~ '\y[0-9]{2}.[0-9]{3}.[0-9]{3}/[0-9]{4}-[0-9]{2}\y';


-- Cleaning not valids CPF and CNPJ
-- if even after formating the cpf/cnpj the value is in invalid format then remove it
delete
from public.file_data fd
where not fd.cpf ~ '\y[0-9]{3}.[0-9]{3}.[0-9]{3}-[0-9]{2}\y'
or not fd.loja_mais_frequente ~ '\y[0-9]{2}.[0-9]{3}.[0-9]{3}/[0-9]{4}-[0-9]{2}\y'
or not fd.loja_ultima_compra ~ '\y[0-9]{2}.[0-9]{3}.[0-9]{3}/[0-9]{4}-[0-9]{2}\y';
