package DriverAction

import (
	"github.com/tebeka/selenium"
	"time"
)

func DriverGet(URL string, wd selenium.WebDriver) {
	err := wd.Get(URL)
	if err != nil {
		panic(err)
	}
	script := `
	    var imgs = document.getElementsByTagName('img');
	    for (var i = 0; i < imgs.length; i++) {
	        imgs[i].parentNode.removeChild(imgs[i]);
	    }
	    var style = document.createElement('style');
	    style.innerHTML = 'img { opacity: 0 }';
	    document.head.appendChild(style);
	    var style = document.createElement('style');
	    style.innerHTML = '* { background-image: none !important; }';
	    document.head.appendChild(style);
	    var style = document.createElement('style');
	    style.innerHTML = '* { color: transparent !important; }';
	    document.head.appendChild(style);
	    var style = document.createElement('style');
	    style.innerHTML = 'img.fade_image { display: none !important; }';
	    document.head.appendChild(style);
	    var style = document.createElement('style');
	    style.innerHTML = '* { transition: paused !important; }';
	    document.head.appendChild(style);
	`
	wd.ExecuteScript(script, nil)
}

func DriverElementWaitAndClick(wd selenium.WebDriver, xpath string) error {
	byXpath := selenium.ByXPATH
	checkTime := 0
	for {
		element, err := wd.FindElement(byXpath, xpath)
		if err == nil {
			isEnabled, err1 := element.IsEnabled()
			if err1 == nil && isEnabled {
				err = element.Click()
				if err != nil {
					continue
				} else {
					return nil
				}
			} else if checkTime > 15 {
				return err
			} else {
				checkTime++
				time.Sleep(500 * time.Millisecond)
				continue
			}
		} else if checkTime > 15 {
			return err
		} else {
			checkTime++
			time.Sleep(500 * time.Millisecond)
			continue
		}

	}
}

func DriverwaitForElement(wd selenium.WebDriver, xpath string) (bool, error) {
	for i := 0; i < 5; i++ {
		_, err := wd.FindElement(selenium.ByXPATH, xpath)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			return true, nil
		}

	}
	return false, nil
}
