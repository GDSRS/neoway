CREATE TABLE IF NOT EXISTS public.file_data (
    cpf text,
    private boolean,
    incompleto boolean,
    data_ultima_compra date,
    ticket_medio money,
    ticket_ultima_compra money,
    loja_mais_frequente text,
    loja_ultima_compra text
);

CREATE OR REPLACE FUNCTION public.format_cpf_cnpj(string_value text) RETURNS text as $$
BEGIN
    IF length(string_value) < 14 THEN
        -- if is less then 14 than must be a cpf
        return format('%s.%s.%s-%s',substring(string_value, 1, 3), substring(string_value, 4, 3), substring(string_value, 7, 3), substring(string_value from 10));
    ELSE
        -- if is equal of greater than 14 must be cnpj
        return format('%s.%s.%s/%s-%s',substring(string_value, 1, 2), substring(string_value, 3, 3), substring(string_value, 7, 3), substring(string_value, 10, 4), substring(string_value from 14));
    END IF;
end;
$$ language plpgsql;