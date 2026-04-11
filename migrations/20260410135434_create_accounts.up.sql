begin;

create table if not exists accounts (
    id bigint generated always as identity primary key,
    document_number text not null unique,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint chk_accounts_document_number_no_spaces check (document_number !~ '\s')
);

create trigger trg_accounts_updated_at
before update on accounts
for each row
execute function set_updated_at();

create index idx_accounts_document_number on accounts (document_number);

commit;
