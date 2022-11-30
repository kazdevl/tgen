package thirdparty

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

type IFSampleClient interface {
	GenrateRandomName() (string, error)
}

type SampleClient struct {
}

func NewSampleClient() *SampleClient {
	return &SampleClient{}
}

func (r *SampleClient) GenrateRandomName() (string, error) {
	return "RandomSample", nil
}
