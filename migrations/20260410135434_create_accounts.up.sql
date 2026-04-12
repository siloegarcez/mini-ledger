begin;

create table if not exists accounts (
    id bigint generated always as identity primary key,
    document_number text not null unique,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    -- No spaces allowed in document_number
    constraint chk_accounts_document_number_no_spaces check (document_number !~ '\s'),
    constraint chk_account_document_number_not_empty check (btrim(document_number) != ''),
    constraint chk_accounts_document_number_max_len check (char_length(document_number) <= 15)
);

create trigger trg_accounts_updated_at
before update on accounts
for each row
execute function set_updated_at();

create index idx_accounts_document_number on accounts (document_number);

commit;
