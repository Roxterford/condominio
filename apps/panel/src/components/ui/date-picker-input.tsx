"use client";

import { Calendar } from "@/components/ui/calendar";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { CalendarIcon } from "lucide-react";
import * as React from "react";

export interface DatePickerInputProps {
  placeholder?: string;
  value?: Date | undefined;
  onChange?: (date: Date | undefined) => void;
  locale?: string;
  id?: string;
  disabled?: boolean;
  required?: boolean;
  className?: string;
}

function formatDate(date: Date | undefined, locale: string = "en-US") {
  if (!date) {
    return "";
  }

  return date.toLocaleDateString(locale, {
    day: "2-digit",
    month: "long",
    year: "numeric",
  });
}

function isValidDate(date: Date | undefined) {
  if (!date) {
    return false;
  }
  return !isNaN(date.getTime());
}

export function DatePickerInput({
  placeholder,
  value,
  onChange,
  locale = "en-US",
  id = "date-picker",
  disabled = false,
  required = false,
  className,
}: DatePickerInputProps) {
  const [open, setOpen] = React.useState(false);
  const [internalDate, setInternalDate] = React.useState<Date | undefined>(
    value,
  );
  const [month, setMonth] = React.useState<Date | undefined>(
    value || new Date(),
  );

  // Use controlled value if provided, otherwise use internal state
  const date = value !== undefined ? value : internalDate;
  const displayValue = formatDate(date, locale);

  return (
    <InputGroup>
      <InputGroupInput
        id={id}
        value={displayValue}
        placeholder={placeholder ?? ""}
        disabled={disabled}
        required={required}
        className={className}
        onChange={(e) => {
          const inputDate = new Date(e.target.value);
          if (isValidDate(inputDate)) {
            if (value === undefined) {
              setInternalDate(inputDate);
            }
            setMonth(inputDate);
            onChange?.(inputDate);
          }
        }}
        onKeyDown={(e) => {
          if (e.key === "ArrowDown") {
            e.preventDefault();
            setOpen(true);
          }
        }}
      />
      <InputGroupAddon align="inline-end">
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <InputGroupButton
              id="date-picker"
              variant="ghost"
              size="icon-xs"
              aria-label="Select date"
            >
              <CalendarIcon />
              <span className="sr-only">Select date</span>
            </InputGroupButton>
          </PopoverTrigger>
          <PopoverContent
            className="w-auto overflow-hidden p-0"
            align="end"
            alignOffset={-8}
            sideOffset={10}
          >
            <Calendar
              mode="single"
              selected={date}
              month={month}
              onMonthChange={setMonth}
              onSelect={(selectedDate) => {
                if (value === undefined) {
                  setInternalDate(selectedDate);
                }
                setMonth(selectedDate);
                onChange?.(selectedDate);
                setOpen(false);
              }}
            />
          </PopoverContent>
        </Popover>
      </InputGroupAddon>
    </InputGroup>
  );
}
