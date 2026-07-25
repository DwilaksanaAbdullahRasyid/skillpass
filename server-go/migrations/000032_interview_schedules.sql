-- server-go/migrations/000032_interview_schedules.sql
-- Interview scheduling for the ATS pipeline: when a company moves a candidate
-- to the "interviewed" stage they record a date/time and place (onsite address
-- or online link). One row per scheduling event so reschedules keep history;
-- the latest row is the active appointment.

CREATE TABLE IF NOT EXISTS interview_schedules (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    scheduled_at   TIMESTAMPTZ NOT NULL,
    mode           TEXT NOT NULL DEFAULT 'onsite', -- 'onsite' | 'online'
    location       TEXT,        -- address (onsite)
    meeting_link   TEXT,        -- meeting URL (online)
    interviewer    TEXT,
    notes          TEXT,
    created_by     UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS interview_schedules_application_idx
    ON interview_schedules(application_id, scheduled_at DESC);
