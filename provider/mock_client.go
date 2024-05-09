// Copyright 2024 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package provider

import (
	"context"

	metal "github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) FindDeviceById(ctx context.Context, id string) metal.ApiFindDeviceByIdRequest {
	args := m.Called(ctx, id)
	return args.Get(0).(metal.ApiFindDeviceByIdRequest)
}

func (m *MockClient) FindProjectDevices(ctx context.Context, id string) metal.ApiFindProjectDevicesRequest {
	args := m.Called(ctx, id)
	return args.Get(0).(metal.ApiFindProjectDevicesRequest)
}

func (m *MockClient) CreateDevice(ctx context.Context, id string) metal.ApiCreateDeviceRequest {
	args := m.Called(ctx, id)
	return args.Get(0).(metal.ApiCreateDeviceRequest)
}

func (m *MockClient) DeleteDevice(ctx context.Context, id string) metal.ApiDeleteDeviceRequest {
	args := m.Called(ctx, id)
	return args.Get(0).(metal.ApiDeleteDeviceRequest)
}

func (m *MockClient) PerformAction(ctx context.Context, id string) metal.ApiPerformActionRequest {
	args := m.Called(ctx, id)
	return args.Get(0).(metal.ApiPerformActionRequest)
}
