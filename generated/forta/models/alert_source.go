// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// AlertSource alert source
//
// swagger:model AlertSource
type AlertSource struct {

	// block
	Block *AlertBlock `json:"block,omitempty"`

	// bot
	Bot *AlertBot `json:"bot,omitempty"`

	// source event
	SourceEvent *AlertSourceEvent `json:"sourceEvent,omitempty"`

	// transaction hash
	// Example: 0x7040dd33cbfd3e9d880da80cb5f3697a717fc329abd0251f3dcd51599ab67b0a
	TransactionHash string `json:"transactionHash,omitempty"`
}

// Validate validates this alert source
func (m *AlertSource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBlock(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBot(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSourceEvent(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertSource) validateBlock(formats strfmt.Registry) error {
	if swag.IsZero(m.Block) { // not required
		return nil
	}

	if m.Block != nil {
		if err := m.Block.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("block")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("block")
			}
			return err
		}
	}

	return nil
}

func (m *AlertSource) validateBot(formats strfmt.Registry) error {
	if swag.IsZero(m.Bot) { // not required
		return nil
	}

	if m.Bot != nil {
		if err := m.Bot.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bot")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("bot")
			}
			return err
		}
	}

	return nil
}

func (m *AlertSource) validateSourceEvent(formats strfmt.Registry) error {
	if swag.IsZero(m.SourceEvent) { // not required
		return nil
	}

	if m.SourceEvent != nil {
		if err := m.SourceEvent.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sourceEvent")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("sourceEvent")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this alert source based on the context it is used
func (m *AlertSource) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBlock(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateBot(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSourceEvent(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertSource) contextValidateBlock(ctx context.Context, formats strfmt.Registry) error {

	if m.Block != nil {

		if swag.IsZero(m.Block) { // not required
			return nil
		}

		if err := m.Block.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("block")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("block")
			}
			return err
		}
	}

	return nil
}

func (m *AlertSource) contextValidateBot(ctx context.Context, formats strfmt.Registry) error {

	if m.Bot != nil {

		if swag.IsZero(m.Bot) { // not required
			return nil
		}

		if err := m.Bot.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bot")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("bot")
			}
			return err
		}
	}

	return nil
}

func (m *AlertSource) contextValidateSourceEvent(ctx context.Context, formats strfmt.Registry) error {

	if m.SourceEvent != nil {

		if swag.IsZero(m.SourceEvent) { // not required
			return nil
		}

		if err := m.SourceEvent.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sourceEvent")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("sourceEvent")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AlertSource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AlertSource) UnmarshalBinary(b []byte) error {
	var res AlertSource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
