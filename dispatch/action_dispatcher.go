package dispatch

import (
	"context"
	"errors"
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/output"
	"github.com/AWoelfel/go-multi-ssh/utils"
)

type dispatchError struct {
	meta output.RemoteMeta
	err  error
}

func (d *dispatchError) Error() string {
	return fmt.Errorf("action failed for %q (%w)", d.meta.Label(), d.err).Error()
}

func (d *dispatchError) Unwrap() error {
	return d.err
}

func Dispatch(ctx context.Context, params []RemoteAction) error {

	stdOutChannel := output.OutputSinkFromContext(ctx, output.StdOutChannel)
	defer stdOutChannel.Wait()

	stdErrChannel := output.OutputSinkFromContext(ctx, output.StdErrChannel)
	defer stdErrChannel.Wait()

	errorChannel := make(chan *dispatchError)
	defer close(errorChannel)

	for i := 0; i < len(params); i++ {
		go func(idx int, actionParam RemoteAction) {
			channelEntry := &dispatchError{
				meta: actionParam.Remote(),
				err:  nil,
			}

			defer func() {
				if rErr := recover(); rErr != nil {
					switch convertedError := rErr.(type) {
					case error:
						channelEntry.err = utils.FromErrors(channelEntry.err, convertedError)
					case string:
						channelEntry.err = utils.FromErrors(channelEntry.err, errors.New(convertedError))
					default:
						channelEntry.err = utils.FromErrors(channelEntry.err, fmt.Errorf("unknown error : %v", convertedError))
					}
				}

				errorChannel <- channelEntry
			}()

			channelEntry.err = actionParam.Run(ctx, stdOutChannel, stdErrChannel)

		}(i, params[i])
	}

	var res []error

	for i := 0; i < len(params); i++ {
		if result := <-errorChannel; result.err != nil {
			stdErrChannel.WriteLine(result.meta, result.err.Error())
			res = append(res, result)
		}
	}

	return utils.FromErrors(res...)
}
