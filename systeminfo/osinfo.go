package systeminfo

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/host"
)

type HardwareInfo struct {
	Hostname        string                  `json:"hostname"`
	Platform        string                  `json:"platform"`
	PlatformVersion string                  `json:"platform_version"`
	VideoChip       []Win32_VideoController `json:"video_chip"`
	NetAdapter      []Win32_NetworkAdapter  `json:"net_adapter"`
	MotherBoard     string                  `json:"mother_board"`
	Memory          []Win32_PhysicalMemory  `json:"memory"`
	HardDisk        []Win32_DiskDrive       `json:"hard_disk"`
	Cameras         string                  `json:"cameras"`
	BIOS            string                  `json:"bios"`
	Processor       []Win32_Processor       `json:"processor"`
}

type Win32_VideoController struct {
	Name        string `json:"name"`
	AdapterRAM  uint32 `json:"adapter_ram"`
	PNPDeviceID string `json:"pnp_device_id"`
}

type Win32_DiskDrive struct {
	PNPDeviceID string `json:"pnp_device_id"`
	Model       string `json:"model"`
	Size        string `json:"size"`
}

type Win32_PhysicalMemory struct {
	Caption      string `json:"caption"`
	Capacity     string `json:"capacity"`
	Manufacturer string `json:"manufacturer"`
}

type Win32_BaseBoard struct {
	Product      string `json:"product"`
	Manufacturer string `json:"manufacturer"`
}

type Win32_Processor struct {
	Name          string `json:"name"`
	NumberOfCores int    `json:"number_of_cores"`
}

type Win32_NetworkAdapter struct {
	Name            string `json:"name"`
	MACAddress      string `json:"mac_address"`
	NetConnectionID string `json:"net_connection_id"`
}

type Win32_BIOS struct {
	Name    string
	Version string
}

func (H *HardwareInfo) Run() {
	Disk := make([]Win32_DiskDrive, 0)
	q := wmi.CreateQuery(&Disk, "")
	if err := wmi.Query(q, &Disk); err == nil {
		H.HardDisk = Disk
	}

	Video := make([]Win32_VideoController, 0)
	q = wmi.CreateQuery(&Video, "")
	if err := wmi.Query(q, &Video); err == nil {
		H.VideoChip = Video
	}

	Memory := make([]Win32_PhysicalMemory, 0)
	q = wmi.CreateQuery(&Memory, "")
	if err := wmi.Query(q, &Memory); err == nil {
		H.Memory = Memory
	}

	BaseBoard := make([]Win32_BaseBoard, 0)
	q = wmi.CreateQuery(&BaseBoard, "")
	if err := wmi.Query(q, &BaseBoard); err == nil {
		if len(BaseBoard) > 0 {
			H.MotherBoard = fmt.Sprintf("%s %s", BaseBoard[0].Product, BaseBoard[0].Manufacturer)
		}
	}

	NetworkAdapter := make([]Win32_NetworkAdapter, 0)
	q = wmi.CreateQuery(&NetworkAdapter, "where MACAddress IS NOT NULL")
	if err := wmi.Query(q, &NetworkAdapter); err == nil {
		H.NetAdapter = NetworkAdapter
	}

	Processor := make([]Win32_Processor, 0)
	q = wmi.CreateQuery(&Processor, "")
	if err := wmi.Query(q, &Processor); err == nil {
		for _, v := range Processor {
			ProcessorCoreNum += v.NumberOfCores
		}
		H.Processor = Processor
	}

	BIOS := make([]Win32_BIOS, 0)
	q = wmi.CreateQuery(&BIOS, "")
	if err := wmi.Query(q, &BIOS); err == nil {
		if len(BIOS) > 0 {
			H.BIOS = fmt.Sprintf("%s %s", BIOS[0].Name, BIOS[0].Version)
		}
	}

	if hInfo, err := host.Info(); err == nil {
		H.Hostname = hInfo.Hostname
		H.Platform = hInfo.Platform
		H.PlatformVersion = hInfo.PlatformVersion
	}
}
