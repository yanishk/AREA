import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import 'globalData.dart';

class CustomButton extends StatelessWidget {
  const CustomButton({
    super.key,
    required this.text,
    required this.onPressed,
    this.backgroundColor = Colors.white,
    this.size = const Size(double.infinity, 70),
    this.padding = const EdgeInsets.symmetric(),
    this.textPaddingLeft = 0,
    this.image,
    this.borderCircularRadius = 32.0,
    this.shadowSize = 3,
    this.isActivated = true,
  });

  final Color backgroundColor;
  final Widget? text;
  final VoidCallback? onPressed;
  final Size size;
  final EdgeInsets padding;
  final double textPaddingLeft;
  final Image? image;
  final double borderCircularRadius;
  final double shadowSize;
  final bool isActivated;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding,
      child: ElevatedButton(
        style: ElevatedButton.styleFrom(
          backgroundColor: backgroundColor,
          elevation: shadowSize,
          minimumSize: size,
          enableFeedback: isActivated,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(borderCircularRadius),
          ),
        ),
        onPressed: onPressed,
        child: Stack(
          alignment: Alignment.center,
          children: [
            if (image != null)
              Positioned(
                left: 0,
                child: image!,
              ),
            Padding(
                padding: EdgeInsets.only(left: textPaddingLeft), child: text),
            const Center(),
          ],
        ),
      ),
    );
  }
}

class CustomClassicForm extends StatelessWidget {
  const CustomClassicForm({
    super.key,
    required this.hint,
    this.controller,
    this.obscureText = false,
    this.padding = const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
    this.enabled = true,
    this.fillCollor = Colors.white,
    this.borderRadius = 10,
    this.borderWidth = 5,
    this.focusedBorderColor = Colors.black,
    this.disabledBorderColor = const Color.fromARGB(255, 238, 238, 238),
    this.enabledBorderColor = const Color.fromARGB(255, 238, 238, 238),
    this.maxLines = 1,
    this.textColor = Colors.black,
  });

  final String hint;
  final TextEditingController? controller;
  final bool obscureText;
  final EdgeInsetsGeometry padding;
  final bool enabled;
  final Color fillCollor;
  final double borderRadius;
  final double borderWidth;
  final Color focusedBorderColor;
  final Color disabledBorderColor;
  final Color enabledBorderColor;
  final int maxLines;
  final Color textColor;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: <Widget>[
        Padding(
          padding: padding,
          child: Container(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(borderRadius),
            ),
            child: TextField(
              maxLines: maxLines,
              minLines: 1,
              style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.w900,
                  color: textColor),
              controller: controller,
              obscureText: obscureText,
              cursorColor: Colors.black,
              decoration: InputDecoration(
                filled: true,
                fillColor: fillCollor,
                labelStyle: const TextStyle(fontSize: 40),
                hintStyle: const TextStyle(
                  fontSize: 20,
                  color: Color.fromARGB(255, 150, 150, 150),
                  fontWeight: FontWeight.bold,
                ),
                hintText: hint,
                enabled: enabled,
                contentPadding:
                    const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(borderRadius),
                  borderSide: BorderSide(
                    color: enabledBorderColor,
                    width: borderWidth,
                  ),
                ),
                disabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(borderRadius),
                  borderSide: BorderSide(
                    color: disabledBorderColor,
                    width: borderWidth,
                  ),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(borderRadius),
                  borderSide: BorderSide(
                    color: focusedBorderColor,
                    width: borderWidth,
                  ),
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class MyDivider extends StatelessWidget {
  const MyDivider({
    super.key,
    this.text = const AutoSizeText(''),
    this.textPadding = const EdgeInsets.symmetric(),
    this.color = Colors.grey,
    this.thickness = 1,
  });

  final AutoSizeText text;
  final EdgeInsets textPadding;
  final Color color;
  final double thickness;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(
          child: Divider(
            thickness: thickness,
            color: color,
          ),
        ),
        Padding(
          padding: textPadding,
          child: text,
        ),
        Expanded(
          child: Divider(
            thickness: thickness,
            color: color,
          ),
        ),
      ],
    );
  }
}

class MyVerticalDivider extends StatelessWidget {
  const MyVerticalDivider({
    super.key,
    this.height,
    this.width,
    this.color = Colors.black,
    this.padding = const EdgeInsets.only(left: 0),
  });

  final double? height;
  final double? width;
  final Color color;
  final EdgeInsetsGeometry padding;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        Padding(
          padding: padding,
          child: SizedBox(
            height: height,
            width: width,
            child: DecoratedBox(
              decoration: BoxDecoration(
                color: color,
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class PlaceableWidget extends StatelessWidget {
  const PlaceableWidget({
    super.key,
    required this.child,
    this.left = 0,
    this.right = 0,
    this.top = 0,
    this.bottom = 0,
  });

  final Widget child;
  final double left;
  final double right;
  final double top;
  final double bottom;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        Padding(
          padding: EdgeInsets.only(
              left: left, right: right, top: top, bottom: bottom),
          child: child,
        ),
      ],
    );
  }
}

class LeftWidget extends StatelessWidget {
  const LeftWidget({
    super.key,
    required this.child,
    this.padding = const EdgeInsets.symmetric(horizontal: 20),
  });

  final Widget child;
  final EdgeInsetsGeometry padding;
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding,
      child: Row(
        children: [
          child,
        ],
      ),
    );
  }
}

class CustomTextButton extends StatelessWidget {
  const CustomTextButton({
    super.key,
    required this.text,
    required this.onPressed,
    this.padding,
  });
  final text;
  final onPressed;
  final padding;

  @override
  Widget build(BuildContext context) {
    return IntrinsicWidth(
      child: TextButton(
        onPressed: onPressed,
        style: TextButton.styleFrom(
          padding: padding,
        ),
        child: text,
      ),
    );
  }
}

class TwoChoiceAlertDialog extends StatelessWidget {
  const TwoChoiceAlertDialog({
    super.key,
    required this.firstChoiceText,
    required this.secondChoiceText,
    required this.firstChoiceonPressed,
    required this.secondChoiceonPressed,
  });

  final Text firstChoiceText;
  final Text secondChoiceText;
  final dynamic firstChoiceonPressed;
  final dynamic secondChoiceonPressed;

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      content: const SingleChildScrollView(
        child: ListBody(
          children: <Widget>[
            Center(
              child: Text("Are you sure?",
                  style: TextStyle(
                    fontSize: 40,
                    fontWeight: FontWeight.w900,
                    color: Colors.black,
                  )),
            )
          ],
        ),
      ),
      actions: <Widget>[
        Center(
          child: CustomButton(
            backgroundColor: const Color.fromARGB(255, 34, 34, 34),
            text: firstChoiceText,
            onPressed: firstChoiceonPressed,
            size: const Size(double.infinity, 55),
            padding: const EdgeInsets.symmetric(horizontal: 15.0),
          ),
        ),
        Center(
          child: CustomTextButton(
            text: secondChoiceText,
            onPressed: secondChoiceonPressed,
          ),
        ),
      ],
    );
  }
}

class WidgetCard extends StatelessWidget {
  const WidgetCard({
    super.key,
    required this.child,
    this.backgroundColor = Colors.black,
    this.borderRadius = 4.0,
  });

  final Widget child;
  final Color backgroundColor;
  final double borderRadius;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return Container(
              decoration: BoxDecoration(
                color: backgroundColor,
                border: Border.all(
                  color: backgroundColor,
                  width: 2.0,
                ),
                borderRadius: BorderRadius.circular(borderRadius),
              ),
              child: Column(
                children: [
                  child,
                ],
              ));
        },
      ),
    );
  }
}

class TextWithBackground extends StatelessWidget {
  const TextWithBackground({
    super.key,
    required this.text,
    this.backgroundColor = Colors.white,
    this.padding = const EdgeInsets.symmetric(horizontal: 15.0),
    this.borderWidth = 1.0,
    this.borderRadius = 15.0,
    this.fontSize = 30,
    this.textColor = Colors.black,
    this.maxLines = 1,
    this.width,
    this.height,
    this.alignment,
  });

  final Color backgroundColor;
  final EdgeInsetsGeometry? padding;
  final double borderWidth;
  final double borderRadius;
  final String text;
  final double fontSize;
  final Color textColor;
  final int maxLines;
  final double ?width;
  final double ?height;
  final Alignment ?alignment;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: padding,
      width: width,
      height: height,
      alignment: alignment,
      decoration: BoxDecoration(
        color: backgroundColor,
        border: Border.all(
          color: backgroundColor,
          width: borderWidth,
        ),
        borderRadius: BorderRadius.circular(borderRadius),
      ),
      child: AutoSizeText(
        text,
        style: TextStyle(
          fontSize: fontSize,
          fontWeight: FontWeight.bold,
          color: textColor,
          backgroundColor: backgroundColor,
        ),
        maxLines: maxLines,
      ),
    );
  }
}

class DoubleWidget extends StatelessWidget {
  const DoubleWidget({
    super.key,
    required this.firstChild,
    required this.secondChild,
    this.padding = const EdgeInsets.only(left: 50),
    this.mainAxisAlignment = MainAxisAlignment.spaceEvenly,
  });

  final Widget firstChild;
  final Widget secondChild;
  final EdgeInsetsGeometry padding;
  final MainAxisAlignment mainAxisAlignment;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: mainAxisAlignment,
      children: [
        firstChild,
        Padding(padding: padding),
        secondChild,
      ],
    );
  }
}

class SafeWebImage extends StatelessWidget {
  const SafeWebImage({
    super.key,
    required this.url,
    this.width = 75,
    this.height = 75,
  });

  final String url;
  final double width;
  final double height;

  @override
  Widget build(BuildContext context) {
    return FadeInImage(
      image: NetworkImage(url),
      placeholder: AssetImage('assets/icon/NoImage.webp'),
      imageErrorBuilder: (context, error, stackTrace) {
        return Image.asset(
          'assets/icon/NoImage.webp',
          width: width,
          height: height,
        );
      },
      width: width,
      height: height,
    );
  }
}

class HexColor extends Color {
  static int _getColorFromHex(String hexColor) {
    hexColor = hexColor.toUpperCase().replaceAll("#", "");
    if (hexColor.length == 6) {
      hexColor = "FF$hexColor";
    }
    return int.parse(hexColor, radix: 16);
  }

  HexColor(final String hexColor) : super(_getColorFromHex(hexColor));
}

class SquareServicesCard extends StatelessWidget {
  const SquareServicesCard({
    super.key,
    required this.name,
    required this.color,
    this.imageUrl = "",
    this.onTap,
  });

  final String name;
  final Color color;
  final String imageUrl;
  final void Function()? onTap;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
      child: Container(
        width: (MediaQuery.sizeOf(Globaldata.myContext).width / 2) - 25,
        height: (MediaQuery.sizeOf(Globaldata.myContext).width / 2) - 25,
        decoration: BoxDecoration(
          color: color,
          border: Border.all(
            color: color,
            width: 2.0,
          ),
          borderRadius: BorderRadius.circular(6.0),
          boxShadow: [
            BoxShadow(
              color: Colors.grey.withOpacity(0.5),
              spreadRadius: 1,
              blurRadius: 8,
              offset: Offset(0, 3),
            ),
          ],
        ),
        child: InkWell(
          onTap: onTap,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SafeWebImage(
                url: "${Globaldata.domainName}$imageUrl",
                width: 50,
                height: 50,
              ),
              AutoSizeText(
                name,
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: Colors.white
                ),
                maxLines: 1,
              )
            ],
          ),
        ),
      ),
    );
  }
}

class EditWheel extends StatelessWidget {
  const EditWheel({
    super.key,
    this.onPressed,
    this.color = Colors.white,
  });
  final void Function()? onPressed;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: onPressed,
      icon: Icon(
        Icons.settings,
        size: 32,
        color: color,
      ),
    );
  }
}
