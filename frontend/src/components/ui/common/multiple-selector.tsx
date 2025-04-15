import * as React from "react";
import { useRef, useState } from "react";
import { Command as CommandPrimitive } from "cmdk";
import { X } from "lucide-react";
import { Button } from "@/components/ui/lib/button";

export interface Option {
  value: string;
  label: string;
}

interface MultipleSelectorProps {
  options?: Option[];
  value?: Option[];
  onValueChange?: (value: Option[]) => void;
  placeholder?: string;
  disabled?: boolean;
  // Новые пропсы
  inputValue?: string;
  onInputChange?: (value: string) => void;
}

export default function MultipleSelector({
                                           options = [],
                                           value = [],
                                           onValueChange,
                                           placeholder = "Search...",
                                           disabled = false,
                                           // Новые пропсы
                                           inputValue = "",
                                           onInputChange
                                         }: MultipleSelectorProps) {
  const [selectedOptions, setSelectedOptions] = useState<Option[]>(value);
  // Убрали внутреннее состояние inputValue
  const [isOpen, setIsOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const filteredOptions = options.filter(option =>
    !selectedOptions.some(selected => selected.value === option.value) &&
    option.label.toLowerCase().includes(inputValue.toLowerCase())
  );

  const handleSelect = (option: Option) => {
    const newValue = [...selectedOptions, option];
    setSelectedOptions(newValue);
    onValueChange?.(newValue);
    onInputChange?.(""); // Сбрасываем input
    inputRef.current?.focus();
  };


  const handleRemove = (option: Option) => {
    const newValue = selectedOptions.filter(p => p.value !== option.value);
    setSelectedOptions(newValue);
    onValueChange?.(newValue);
  };

  const handleClearAll = () => {
    setSelectedOptions([]);
    onValueChange?.([]);
    inputRef.current?.focus();
  };

  return (
    <div className="relative w-full space-y-2" data-disabled={disabled}>
      <div className="flex flex-wrap gap-2">
        {selectedOptions.map((opt) => (
          <div
            key={opt.value}
            className="flex items-center gap-1 rounded-full bg-primary text-primary-foreground px-3 py-1 text-sm"
          >
            <span>{opt.label}</span>
            <button
              type="button"

              className="rounded-full hover:bg-primary/80 ml-1"
              disabled={disabled}
            >
              <X className="h-4 w-4" onClick={() => handleRemove(opt)} />
            </button>
          </div>
        ))}
      </div>

      <CommandPrimitive className="relative">
        <div
          className="group rounded-md border border-input px-3 py-2 text-sm ring-offset-background focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-2">
          <CommandPrimitive.Input
            ref={inputRef}
            placeholder={placeholder}
            value={inputValue} // Управляемое значение
            onValueChange={onInputChange} // Пробрасываем изменения
            onFocus={() => !disabled && setIsOpen(true)}
            onBlur={() => setTimeout(() => setIsOpen(false), 100)}
            className="w-full bg-transparent outline-none placeholder:text-muted-foreground"
            disabled={disabled}
          />
        </div>

        {isOpen && !disabled && (
          <CommandPrimitive.List
            className="absolute top-full w-full mt-2 rounded-md border bg-popover text-popover-foreground shadow-md animate-in fade-in-0 zoom-in-95"
            onMouseDown={(e) => e.preventDefault()}
          >
            <CommandPrimitive.Empty className="px-3 py-1.5 text-sm text-muted-foreground">
              No results found
            </CommandPrimitive.Empty>

            {filteredOptions.map((opt) => (
              <CommandPrimitive.Item
                key={opt.value}
                onSelect={() => handleSelect(opt)}
                className="cursor-default px-3 py-1.5 text-sm aria-selected:bg-accent aria-selected:text-accent-foreground hover:bg-accent hover:text-accent-foreground"
              >
                {opt.label}
              </CommandPrimitive.Item>
            ))}
          </CommandPrimitive.List>
        )}
      </CommandPrimitive>

      {selectedOptions.length > 0 && !disabled && (
        <Button
          variant="ghost"
          size="sm"
          onClick={handleClearAll}
          className="h-8 text-destructive hover:text-destructive/80"
          type="button"
        >
          Clear All
        </Button>
      )}
    </div>
  );
}